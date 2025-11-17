# ISO8583 Message Pack/Unpack/Marshal/Unmarshal Specification

This document specifies the behavior of the core message packing, unpacking, marshaling, and unmarshaling functionality in the ISO8583 library. These implementations have been removed to enable spec-driven reimplementation.

## Overview

The Message type provides four core operations for working with ISO8583 messages:
- **Pack**: Converts a Message into its binary wire format
- **Unpack**: Parses binary wire format into a Message
- **Marshal**: Populates Message fields from Go structs
- **Unmarshal**: Extracts Message field values into Go structs

## Data Structures

### Message Structure
```go
type Message struct {
    spec         *MessageSpec
    cachedBitmap *field.Bitmap
    fields       map[int]field.Field  // All fields according to spec
    mu           sync.Mutex            // Guards fieldsMap
    fieldsMap    map[int]struct{}     // Tracks which fields were set
}
```

### Field Indices
- Index 0: MTI (Message Type Indicator)
- Index 1: Bitmap
- Index 2+: Data fields

## Pack Operation

### Purpose
Converts a Message into its binary wire format for transmission over the network.

### Method Signatures
```go
func (m *Message) Pack() ([]byte, error)
func (m *Message) wrapErrorPack() ([]byte, error)  // internal
func (m *Message) pack() ([]byte, error)           // internal
```

### Algorithm

1. **Thread Safety**: Lock the message mutex for the duration of the operation

2. **Initialize Bitmap**: Reset the bitmap to clear any previous state

3. **Determine Packable Fields**:
   - Get all field IDs from `fieldsMap` (fields that have been set)
   - Always include bitmap (index 1)
   - Sort field IDs in ascending order
   - Return as `packableFieldIDs()`

4. **Set Bitmap Bits**:
   - Iterate through packable field IDs
   - Skip MTI (index 0) and bitmap (index 1)
   - Skip bitmap presence indicator bits (e.g., 1, 65, 129, 193 for default 64-bit bitmap)
   - For each remaining field ID, set the corresponding bit in the bitmap

5. **Pack Fields in Order**:
   - Iterate through sorted packable field IDs
   - Skip bitmap presence indicator bits (except the first bitmap itself at index 1)
   - For each field:
     - Look up the field in the `fields` map
     - Call `field.Pack()` to get the packed bytes
     - Append packed bytes to the output buffer
   - If a field doesn't exist in the spec, return an error

6. **Error Handling**:
   - Wrap any errors in `*iso8583errors.PackError`
   - Include field ID and description in error messages

### Bitmap Presence Bits
- For a 64-bit bitmap: bits 1, 65, 129, 193, etc. indicate presence of additional bitmaps
- These bits are NOT packed as regular fields
- They are automatically managed by the bitmap itself during Pack/Unpack

### Expected Behavior
- Returns packed bytes ready for network transmission
- MTI is always first in the output
- Bitmap is always second in the output
- Fields are packed in ascending order by field ID
- Empty/unset fields are not included in the output
- Thread-safe operation

### Error Conditions
- Field specified in fieldsMap but not in spec: "failed to pack field %d: no specification found"
- Field Pack() fails: "failed to pack field %d (%s): %w" (includes field description)
- General pack failure: "failed to pack message: %w"

## Unpack Operation

### Purpose
Parses binary wire format into a Message, populating fields according to the spec.

### Method Signatures
```go
func (m *Message) Unpack(src []byte) error
func (m *Message) wrapErrorUnpack(src []byte) error  // internal
func (m *Message) unpack(src []byte) (string, error) // internal, returns fieldID on error
```

### Algorithm

1. **Thread Safety**: Lock the message mutex for the duration of the operation

2. **Initialize State**:
   - Reset `fieldsMap` to empty map
   - Reset the bitmap
   - Initialize offset to 0

3. **Unpack MTI** (index 0):
   - Call `fields[mtiIdx].Unpack(src)` 
   - Mark MTI as set in `fieldsMap`
   - Advance offset by bytes read
   - On error: return "0" as field ID

4. **Unpack Bitmap** (index 1):
   - Call `fields[bitmapIdx].Unpack(src[off:])`
   - Bitmap automatically sets itself in `fieldsMap`
   - Advance offset by bytes read
   - On error: return "1" as field ID

5. **Unpack Data Fields**:
   - Iterate from field 2 to `bitmap.Len()` (maximum field number)
   - Skip bitmap presence indicator bits (e.g., 65, 129, 193)
   - For each field ID where `bitmap.IsSet(id)` returns true:
     - Look up field in `fields` map
     - If field doesn't exist in spec, return error with field ID
     - Call `field.Unpack(src[off:])`
     - Mark field as set in `fieldsMap`
     - Advance offset by bytes read
     - On error: return field ID as string

6. **Error Handling**:
   - Wrap errors in `*iso8583errors.UnpackError`
   - Include field ID, error, and raw message in UnpackError
   - Return empty string as field ID on success

### Expected Behavior
- Parses complete message from wire format
- Only unpacks fields indicated by bitmap
- Maintains field order during unpacking
- Thread-safe operation
- Clears previous state before unpacking

### Error Conditions
- MTI unpack fails: "failed to unpack MTI: %w"
- Bitmap unpack fails: "failed to unpack bitmap: %w"
- Field not in spec: "failed to unpack field %d: no specification found"
- Field unpack fails: "failed to unpack field %d (%s): %w"

## Marshal Operation

### Purpose
Populates Message fields from a Go struct, mapping struct fields to ISO8583 fields using tags.

### Method Signatures
```go
func (m *Message) Marshal(v interface{}) error
func (m *Message) marshalStruct(dataStruct reflect.Value) error  // internal
```

### Algorithm

1. **Thread Safety**: Lock the message mutex for the duration of the operation

2. **Validate Input**:
   - Return nil if v is nil
   - Get reflect.Value of v
   - If pointer or interface, dereference to get underlying value
   - Verify underlying value is a struct, else return error

3. **Iterate Struct Fields**:
   - For each field in the struct:
     - Parse the `iso8583` or `index` tag to get field ID
     - If field has valid index tag (ID >= 0):
       - Get corresponding message field from spec
       - If message field doesn't exist, return error
       - Check if struct field is zero value
       - If zero and `keepzero` tag not set, skip this field
       - Call `messageField.Marshal(dataField.Interface())`
       - Mark field as set in `fieldsMap`
     - If field is anonymous embedded struct without index tag:
       - Recursively call `marshalStruct()` on the embedded struct
       - This allows composition of message structs

4. **Tag Format**:
   - `iso8583:"0"` - Field 0 (MTI)
   - `iso8583:"2"` - Field 2
   - `iso8583:"48,keepzero"` - Field 48, include even if zero value
   - `index:"1"` - Legacy format, same as iso8583

5. **Embedded Struct Handling**:
   - Anonymous embedded structs without index tags are traversed recursively
   - Handles pointer and interface types by dereferencing
   - Skips nil embedded structs
   - Non-anonymous fields without tags are ignored

### Expected Behavior
- Maps Go struct fields to ISO8583 message fields
- Supports nested/composite fields
- Supports embedded anonymous structs
- Skips zero-value fields unless `keepzero` is specified
- Thread-safe operation

### Error Conditions
- v is not a struct: "data is not a struct"
- Field ID not in spec: "no message field defined by spec with index: %d"
- Field marshal fails: "failed to set value to field %d: %w"

## Unmarshal Operation

### Purpose
Extracts Message field values into a Go struct, mapping ISO8583 fields to struct fields using tags.

### Method Signatures
```go
func (m *Message) Unmarshal(v interface{}) error
func (m *Message) unmarshalStruct(dataStruct reflect.Value) error  // internal
```

### Algorithm

1. **Thread Safety**: Lock the message mutex for the duration of the operation

2. **Validate Input**:
   - Verify v is a pointer and not nil, else return error
   - Dereference pointer to get struct
   - Verify underlying value is a struct, else return error

3. **Iterate Struct Fields**:
   - For each field in the struct:
     - Parse the `iso8583` or `index` tag to get field ID
     - If field has valid index tag (ID >= 0):
       - Get corresponding message field
       - If message field is nil, skip (field not in spec)
       - If field not set in `fieldsMap`, skip (field not present in message)
       - Handle different struct field kinds:
         - **Pointer/Interface**: Initialize if nil, call `messageField.Unmarshal(dataField.Interface())`
         - **Slice**: Pass reflect.Value so slice can be modified, call `messageField.Unmarshal(dataField)`
         - **Native types**: Call `messageField.Unmarshal(dataField)`
     - If field is anonymous embedded struct without index tag:
       - Recursively call `unmarshalStruct()` on the embedded struct
       - Initialize nil pointers if possible
       - Skip nil embedded structs that can't be initialized

4. **Embedded Struct Handling**:
   - Anonymous embedded structs without index tags are traversed recursively
   - Attempts to initialize nil pointer embedded structs
   - Handles pointer and interface types by dereferencing
   - Skips fields that can't be initialized

### Expected Behavior
- Maps ISO8583 message fields to Go struct fields
- Only unmarshals fields that are set in the message
- Supports nested/composite fields
- Supports embedded anonymous structs
- Initializes nil pointers when possible
- Thread-safe operation

### Error Conditions
- v is not a pointer or is nil: "data is not a pointer or nil"
- v does not point to a struct: "data is not a struct"
- Field unmarshal fails: "failed to get value from field %d: %w"

## Helper Methods

### packableFieldIDs()
```go
func (m *Message) packableFieldIDs() ([]int, error)
```
- Returns sorted list of field IDs to pack
- Always includes bitmap (index 1)
- Includes all fields in `fieldsMap` except bitmap
- Returns fields in ascending order

### MarshalJSON()
```go
func (m *Message) MarshalJSON() ([]byte, error)
```
- Calls `wrapErrorPack()` to generate bitmap and validate message
- Converts field map to string-keyed map
- Uses `field.OrderedMap()` to maintain field order
- Returns JSON bytes

### UnmarshalJSON()
```go
func (m *Message) UnmarshalJSON(b []byte) error
```
- Parses JSON into map of field ID to raw JSON
- For each field:
  - Converts string ID to int
  - Gets field from spec
  - Unmarshals raw JSON into field
  - Marks field as set in `fieldsMap`

### Clone()
```go
func (m *Message) Clone() (*Message, error)
```
- Creates new message with same spec
- Packs current message to bytes
- Unpacks bytes into new message
- Returns cloned message

## Concurrency

All public methods (Pack, Unpack, Marshal, Unmarshal) are thread-safe:
- Lock `m.mu` at the start
- Defer unlock
- Internal methods assume mutex is already locked

## Test Coverage

The implementations have the following test coverage:
- Pack: 100%
- Unpack: 100%
- pack (internal): 90.0%
- unpack (internal): 92.0%
- Marshal: 90.0%
- Unmarshal: 88.9%
- marshalStruct: 91.7%
- unmarshalStruct: 88.9%

## Key Test Scenarios

1. **Basic Pack/Unpack**: MTI + simple fields
2. **Bitmap Presence Bits**: Fields 65, 129, 193 are not packed as data
3. **Three Bitmaps**: Messages with fields beyond 128
4. **Composite Fields**: Nested subfields with positional encoding
5. **Marshal/Unmarshal**: Struct to message and back
6. **Embedded Structs**: Anonymous embedded structs in Marshal/Unmarshal
7. **Zero Values**: keepzero tag behavior
8. **Concurrent Access**: No data races when accessing fields concurrently
9. **Error Handling**: Proper error wrapping and field ID reporting
10. **JSON Encoding**: MarshalJSON/UnmarshalJSON round-trip

## Implementation Notes

1. **Bitmap Management**: The bitmap is automatically managed during Pack/Unpack
2. **Field Order**: Fields must be packed/unpacked in ascending order by ID
3. **Presence Indicators**: Bitmap presence bits (1, 65, 129, etc.) are special-cased
4. **Error Context**: Errors include field ID and description for debugging
5. **State Management**: fieldsMap tracks which fields are set
6. **Reflection**: Marshal/Unmarshal use reflection to map between structs and fields
