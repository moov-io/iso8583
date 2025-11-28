package iso8583

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

	iso8583errors "github.com/moov-io/iso8583/errors"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/utils"
)

var _ json.Marshaler = (*Message)(nil)
var _ json.Unmarshaler = (*Message)(nil)

const (
	mtiIdx    = 0
	bitmapIdx = 1
)

type Message struct {
	spec         *MessageSpec
	cachedBitmap *field.Bitmap

	// to guard fields
	mu sync.Mutex

	// stores all fields according to the spec
	fields map[int]field.Field
}

func NewMessage(spec *MessageSpec) *Message {
	return &Message{
		fields: make(map[int]field.Field),
		spec:   spec,
	}
}

// Deprecated. Use Marshal instead.
func (m *Message) SetData(data interface{}) error {
	return m.Marshal(data)
}

func (m *Message) Bitmap() *field.Bitmap {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.bitmap()
}

// bitmap creates and returns the bitmap field, it's not thread safe
// and should be called from a thread safe function
func (m *Message) bitmap() *field.Bitmap {
	if m.cachedBitmap != nil {
		return m.cachedBitmap
	}

	bitmap, err := m.createField(bitmapIdx)
	if err != nil {
		panic(fmt.Sprintf("required bitmap field is missing: %v", err))
	}

	var ok bool
	m.cachedBitmap, ok = bitmap.(*field.Bitmap)
	if !ok {
		panic("bitmap field is not of type *field.Bitmap")
	}
	m.cachedBitmap.Reset()

	return m.cachedBitmap
}

func (m *Message) MTI(val string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mti, err := m.getOrCreateField(mtiIdx)
	if err != nil {
		panic(fmt.Sprintf("required MTI field is missing: %v", err))
	}
	mti.SetBytes([]byte(val))
}

func (m *Message) GetSpec() *MessageSpec {
	return m.spec
}

func (m *Message) Field(id int, val string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	field, err := m.getOrCreateField(id)
	if err != nil {
		return fmt.Errorf("getting or creating field %d: %w", id, err)
	}

	err = field.SetBytes([]byte(val))
	if err != nil {
		return fmt.Errorf("setting bytes for field %d: %w", id, err)
	}

	return nil
}

func (m *Message) BinaryField(id int, val []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	field, err := m.getOrCreateField(id)
	if err != nil {
		return fmt.Errorf("getting or creating field %d: %w", id, err)
	}

	err = field.SetBytes(val)
	if err != nil {
		return fmt.Errorf("setting bytes for field %d: %w", id, err)
	}

	return nil
}

func (m *Message) GetMTI() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mti, err := m.getOrCreateField(mtiIdx)
	if err != nil {
		return "", fmt.Errorf("getting or creating MTI field: %w", err)
	}

	return mti.String()
}

func (m *Message) GetString(id int) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	field, err := m.getOrCreateField(id)
	if err != nil {
		return "", fmt.Errorf("getting or creating field %d: %w", id, err)
	}

	str, err := field.String()
	if err != nil {
		return "", fmt.Errorf("getting string for field %d: %w", id, err)
	}

	return str, nil
}

func (m *Message) GetBytes(id int) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	field, err := m.getOrCreateField(id)
	if err != nil {
		return nil, fmt.Errorf("getting or creating field %d: %w", id, err)
	}

	bytes, err := field.Bytes()
	if err != nil {
		return nil, fmt.Errorf("getting bytes for field %d: %w", id, err)
	}

	return bytes, nil
}

// GetField returns the field with the given ID. If the field does not exist
// in the message, it will be created based on the message specification. If
// the field ID is not defined in the specification, nil will be returned.
func (m *Message) GetField(id int) field.Field {
	m.mu.Lock()
	defer m.mu.Unlock()

	field, _ := m.getOrCreateField(id)
	return field
}

// Fields returns the copy of the map of the set fields in the message. Be aware
// that fields are live references, so modifying them will affect the message.
func (m *Message) GetFields() map[int]field.Field {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.getFields()
}

// getFields returns the map of the set fields. It assumes that the mutex
// is already locked by the caller.
func (m *Message) getFields() map[int]field.Field {
	return maps.Clone(m.fields)
}

// Pack locks the message, packs its fields, and then unlocks it.
// If any errors are encountered during packing, they will be wrapped
// in a *PackError before being returned.
func (m *Message) Pack() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.wrapErrorPack()

}

// wrapErrorPack calls the core packing logic and wraps any errors in a
// *PackError. It assumes that the mutex is already locked by the caller.
func (m *Message) wrapErrorPack() ([]byte, error) {
	data, err := m.pack()
	if err != nil {
		return nil, &iso8583errors.PackError{Err: err}
	}

	return data, nil
}

// pack contains the core logic for packing the message. This method does not
// handle locking or error wrapping and should typically be used internally
// after ensuring concurrency safety.
func (m *Message) pack() ([]byte, error) {
	packed := []byte{}
	m.bitmap().Reset()

	ids, err := m.packableFieldIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to pack message: %w", err)
	}

	for _, id := range ids {
		// indexes 0 and 1 are for mti and bitmap
		// regular field number startd from index 2
		// do not pack presence bits as well
		if id < 2 || m.bitmap().IsBitmapPresenceBit(id) {
			continue
		}
		m.bitmap().Set(id)
	}

	// pack fields
	for _, i := range ids {
		// do not pack presence bits other than the first one as it's the bitmap itself
		if i != 1 && m.bitmap().IsBitmapPresenceBit(i) {
			continue
		}

		// m.fields[i] must have the field as we got i from packableFieldIDs()
		field := m.fields[i]

		packedField, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack field %d (%s): %w", i, field.Spec().Description, err)
		}
		packed = append(packed, packedField...)
	}

	return packed, nil
}

// Unpack unpacks the message from the given byte slice or returns an error
// which is of type *UnpackError and contains the raw message
func (m *Message) Unpack(src []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.wrapErrorUnpack(src)
}

// wrapErrorUnpack calls the core unpacking logic and wraps any
// errors in a *UnpackError. It assumes that the mutex is already
// locked by the caller.
func (m *Message) wrapErrorUnpack(src []byte) error {
	if fieldID, err := m.unpack(src); err != nil {
		return &iso8583errors.UnpackError{
			Err:        err,
			FieldID:    fieldID,
			RawMessage: src,
		}
	}
	return nil
}

// unpack contains the core logic for unpacking the message. This method does
// not handle locking or error wrapping and should typically be used internally
// after ensuring concurrency safety.
func (m *Message) unpack(src []byte) (string, error) {
	var off int

	m.fields = make(map[int]field.Field)
	m.cachedBitmap = nil

	// This method implicitly also sets m.fields[bitmapIdx]
	bitmap := m.bitmap()

	mti, err := m.createField(mtiIdx)
	if err != nil {
		return strconv.Itoa(mtiIdx), fmt.Errorf("getting or creating MTI field: %w", err)
	}

	read, err := mti.Unpack(src)
	if err != nil {
		return strconv.Itoa(mtiIdx), fmt.Errorf("failed to unpack MTI: %w", err)
	}

	off = read

	// unpack Bitmap
	read, err = bitmap.Unpack(src[off:])
	if err != nil {
		return strconv.Itoa(bitmapIdx), fmt.Errorf("failed to unpack bitmap: %w", err)
	}

	off += read

	for i := 2; i <= m.bitmap().Len(); i++ {
		// skip bitmap presence bits (for default bitmap length of 64 these are bits 1, 65, 129, 193, etc.)
		if m.bitmap().IsBitmapPresenceBit(i) {
			continue
		}

		if m.bitmap().IsSet(i) {
			fl, err := m.getOrCreateField(i)
			if err != nil {
				return strconv.Itoa(i), fmt.Errorf("getting or creating field %d: %w", i, err)
			}

			read, err = fl.Unpack(src[off:])
			if err != nil {
				return strconv.Itoa(i), fmt.Errorf("failed to unpack field %d (%s): %w", i, fl.Spec().Description, err)
			}

			fmt.Printf("Unpacked field %d: %+v\n", i, fl)

			off += read
		}
	}

	return "", nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// by packing the message we will generate bitmap
	// create HEX representation
	// and validate message against the spec
	if _, err := m.wrapErrorPack(); err != nil {
		return nil, err
	}

	fieldMap := m.getFields()
	strFieldMap := map[string]field.Field{}
	for id, field := range fieldMap {
		strFieldMap[strconv.Itoa(id)] = field
	}

	// get only fields that were set
	bytes, err := json.Marshal(field.OrderedMap(strFieldMap))
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal map to bytes")
	}
	return bytes, nil
}

func (m *Message) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var data map[string]json.RawMessage
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for id, rawMsg := range data {
		i, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field %v: could not convert to int", i)
		}

		field, err := m.getOrCreateField(i)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field %d: %w", i, err)
		}

		if err := json.Unmarshal(rawMsg, field); err != nil {
			return utils.NewSafeErrorf(err, "failed to unmarshal field %v", id)
		}
	}

	return nil
}

func (m *Message) packableFieldIDs() ([]int, error) {
	return slices.Sorted(maps.Keys(m.fields)), nil
}

// Clone clones the message by creating a new message from the binary
// representation of the original message
func (m *Message) Clone() (*Message, error) {
	newMessage := NewMessage(m.spec)

	m.mu.Lock()
	bytes, err := m.wrapErrorPack()
	if err != nil {
		m.mu.Unlock()
		return nil, err
	}
	m.mu.Unlock()

	mti, err := m.GetMTI()
	if err != nil {
		return nil, err
	}

	newMessage.MTI(mti)
	newMessage.Unpack(bytes)

	_, err = newMessage.Pack()
	if err != nil {
		return nil, err
	}

	return newMessage, nil
}

// Marshal populates message fields with v struct field values. It traverses
// through the message fields and calls Unmarshal(...) on them setting the v If
// v is not a struct or not a pointer to struct then it returns error.
func (m *Message) Marshal(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v == nil {
		return nil
	}

	dataStruct := reflect.ValueOf(v)

	if dataStruct.Kind() == reflect.Ptr || dataStruct.Kind() == reflect.Interface {
		dataStruct = dataStruct.Elem()
	}

	if dataStruct.Kind() != reflect.Struct {
		return errors.New("data is not a struct")
	}

	err := m.marshalStruct(dataStruct)
	if err != nil {
		return fmt.Errorf("marshaling struct: %w", err)
	}

	return nil
}

// marshalStruct is a helper method that handles the core logic of marshaling a struct.
// It supports anonymous embedded structs by recursively traversing into them when they
// don't have index tags themselves.
func (m *Message) marshalStruct(dataStruct reflect.Value) error {
	// iterate over struct fields
	for i := 0; i < dataStruct.NumField(); i++ {
		structField := dataStruct.Type().Field(i)
		indexTag := field.NewIndexTag(structField)

		// If the field has an index tag, process it normally
		if indexTag.ID >= 0 {
			dataField := dataStruct.Field(i)
			// for non pointer fields we need to check if they are zero
			// and we want to skip them (as specified in the field tag)
			if dataField.IsZero() && !indexTag.KeepZero {
				continue
			}

			messageField, err := m.getOrCreateField(indexTag.ID)
			if err != nil {
				return fmt.Errorf("getting or creating field %d: %w", indexTag.ID, err)
			}

			if err := messageField.Marshal(dataField.Interface()); err != nil {
				return fmt.Errorf("failed to set value to field %d: %w", indexTag.ID, err)
			}

			continue
		}

		// If it's an anonymous embedded struct without an index tag, traverse into it
		if structField.Anonymous {
			fieldValue := dataStruct.Field(i)

			// Handle pointer and interface types
			for fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					break // skip nil embedded structs
				}
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Struct && fieldValue.IsValid() {
				// Recursively process the embedded struct
				if err := m.marshalStruct(fieldValue); err != nil {
					return err
				}
			}
		}
		// Otherwise, skip the field (existing behavior for non-anonymous fields without index tags)
	}

	return nil
}

// Unmarshal populates v struct fields with message field values. It traverses
// through the message fields and calls Unmarshal(...) on them setting the v If
// v  is nil or not a pointer it returns error.
func (m *Message) Unmarshal(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("data is not a pointer or nil")
	}

	// get the struct from the pointer
	dataStruct := rv.Elem()

	if dataStruct.Kind() != reflect.Struct {
		return errors.New("data is not a struct")
	}

	return m.unmarshalStruct(dataStruct)
}

// unmarshalStruct is a helper method that handles the core logic of unmarshaling into a struct.
// It supports anonymous embedded structs by recursively traversing into them when they
// don't have index tags themselves.
func (m *Message) unmarshalStruct(dataStruct reflect.Value) error {
	// iterate over struct fields
	for i := range dataStruct.NumField() {
		structField := dataStruct.Type().Field(i)
		indexTag := field.NewIndexTag(structField)

		// If the field has an index tag, process it normally
		if indexTag.ID >= 0 {
			// skip if field is not set in the message
			messageField := m.fields[indexTag.ID]
			if messageField == nil {
				continue
			}

			dataField := dataStruct.Field(i)
			switch dataField.Kind() { //nolint:exhaustive
			case reflect.Pointer, reflect.Interface:
				if dataField.IsNil() {
					dataField.Set(reflect.New(dataField.Type().Elem()))
				}
				err := messageField.Unmarshal(dataField.Interface())
				if err != nil {
					return fmt.Errorf("failed to get value from field %d: %w", indexTag.ID, err)
				}
			case reflect.Slice:
				// Pass reflect.Value for slices so they can be modified
				err := messageField.Unmarshal(dataField)
				if err != nil {
					return fmt.Errorf("failed to get value from field %d: %w", indexTag.ID, err)
				}
			default: // Native types
				err := messageField.Unmarshal(dataField)
				if err != nil {
					return fmt.Errorf("failed to get value from field %d: %w", indexTag.ID, err)
				}
			}
			continue
		}

		// If it's an anonymous embedded struct without an index tag, traverse into it
		if structField.Anonymous {
			fieldValue := dataStruct.Field(i)

			// Handle pointer and interface types
			for fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					// Try to initialize if possible
					if fieldValue.CanSet() && fieldValue.Kind() == reflect.Ptr {
						fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
						fieldValue = fieldValue.Elem()
					} else {
						break // skip nil embedded structs that can't be initialized
					}
				} else {
					fieldValue = fieldValue.Elem()
				}
			}

			if fieldValue.Kind() == reflect.Struct && fieldValue.IsValid() {
				// Recursively process the embedded struct
				if err := m.unmarshalStruct(fieldValue); err != nil {
					return err
				}
			}
		}
		// Otherwise, skip the field (existing behavior for non-anonymous fields without index tags)
	}

	return nil
}

// UnsetField marks the field with the given ID as not set and replaces it with
// a new zero-valued field. This effectively removes the field's value and excludes
// it from operations like Pack() or Marshal().
func (m *Message) UnsetField(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
}

func (m *Message) unsetField(id int) {
	delete(m.fields, id)
}

func (m *Message) getOrCreateField(id int) (field.Field, error) {
	f := m.fields[id]
	if f == nil {
		return m.createField(id)
	}

	return f, nil
}

func (m *Message) createField(id int) (field.Field, error) {
	specField, ok := m.GetSpec().Fields[id]
	if !ok {
		return nil, fmt.Errorf("failed to create field %d as it does not exist in the spec", id)
	}
	f := field.NewInstanceOf(specField)
	m.fields[id] = f

	return f, nil
}

// UnsetFields marks multiple fields identified by their paths as not set and
// replaces them with new zero-valued fields. Each path should be in the format
// "a.b.c". This effectively removes the fields' values and excludes them from
// operations like Pack() or Marshal().
// Deprecated: use UnsetPath instead.
func (m *Message) UnsetFields(idPaths ...string) error {
	return m.UnsetPath(idPaths...)
}

// UnsetPath marks multiple fields identified by their paths as not set and
// replaces them with new zero-valued fields. Each path should be in the format
// "a.b.c". This effectively removes the fields' values and excludes them from
// operations like Pack() or Marshal().
func (m *Message) UnsetPath(idPaths ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, idPath := range idPaths {
		if idPath == "" {
			continue
		}

		id, path, hasSubpath := strings.Cut(idPath, ".")
		idx, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("conversion of %s to int failed: %w", id, err)
		}

		f := m.fields[idx]
		// field is not set, continue
		if f == nil {
			continue
		}

		// If there's no subpath, unset the entire field
		if !hasSubpath {
			m.unsetField(idx)
			continue
		}

		// Handle composite field with subpaths
		pathUnsetter, ok := f.(field.PathUnsetter)
		if !ok {
			return fmt.Errorf("field %d is not a composite field and its subfields %s cannot be unset", idx, path)
		}

		if err := pathUnsetter.UnsetPath(path); err != nil {
			return fmt.Errorf("failed to unset %s in composite field %d: %w", path, idx, err)
		}
	}

	return nil
}

func (m *Message) MarshalPath(path string, value any) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id, subPath, hasSubPath := strings.Cut(path, ".")
	idx, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("conversion of %s to int failed: %w", id, err)
	}

	f, err := m.getOrCreateField(idx)
	if err != nil {
		return fmt.Errorf("field %d does not exist", idx)
	}

	if hasSubPath {
		mField, ok := f.(field.PathMarshaler)
		if !ok {
			return fmt.Errorf("field %s is not a PathMarshaler", id)
		}

		err := mField.MarshalPath(subPath, value)
		if err != nil {
			return fmt.Errorf("marshaling filed %s: %w", id, err)
		}

		return nil
	}

	err = f.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling field %s: %w", id, err)
	}

	return nil
}

func (m *Message) UnmarshalPath(path string, value any) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id, subPath, hasSubPath := strings.Cut(path, ".")
	idx, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("conversion of %s to int failed: %w", id, err)
	}

	f := m.fields[idx]
	if f == nil {
		// check if field exists in spec
		_, ok := m.spec.Fields[idx]
		if !ok {
			return fmt.Errorf("field %d is not defined in the spec", idx)
		}

		return nil
	}

	if hasSubPath {
		uField, ok := f.(field.PathUnmarshaler)
		if !ok {
			return fmt.Errorf("field %s is not a PathUnmarshaler", id)
		}

		err := uField.UnmarshalPath(subPath, value)
		if err != nil {
			return fmt.Errorf("unmarshaling filed %s: %w", id, err)
		}

		return nil
	}

	err = f.Unmarshal(value)
	if err != nil {
		return fmt.Errorf("unmarshaling field %s: %w", id, err)
	}

	return nil
}
