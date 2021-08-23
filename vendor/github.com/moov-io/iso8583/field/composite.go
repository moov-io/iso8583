package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/moov-io/iso8583/padding"
)

var _ Field = (*Composite)(nil)

// Composite is a wrapper object designed to hold ISO8583 subfields and
// subelements.  Because Composite handles both of these usecases generically,
// we refer to them collectively as 'subfields' throughout the receiver's
// documentation and error messages.
//
// Composite handles aggregate fields of the following format:
// - Length (if variable)
// - []Subfield
//
// Where the subfield structure is assumed to be as follows:
// - Subfield ID (if Composite.Spec().IDLength > 0)
// - Subfield Length (if variable)
// - Subfield data
//
// Composite behaves in two modes depending on whether subfield IDs need to be
// explicitly handled or not. This is configured by setting Spec.IDLength.
//
// When populated, Composite handles the packing and unpacking subfield IDs on
// their behalf. However, responsibility for packing and unpacking both the
// length and data of subfields is delegated to the subfields themselves.
// Therefore, their specs should be configured to handle such behavior.
//
// If Spec.IDLength > 0, Composite leverages Spec.Enc to unpack subfields
// regardless of order based on their IDs. Similarly, it is also able to handle
// non-present subfields by relying on the existence of their IDs.
//
// If Spec.IDLength == 0, Composite only unpacks subfields ordered by ID. The absence
// of IDs in the data means that the receiver is not able to handle non-present
// subfields either.
//
// For the sake of determinism, packing of subfields is executed in order of ID
// regardless of the value of Spec.IDLength.
//
// Padding is not supported by Composite. Responsibility for this is delegated
// recursively to the subfields themselves.
type Composite struct {
	spec *Spec

	orderedSpecFieldIDs []int
	idToFieldMap        map[int]Field

	fieldsMap map[int]struct{}
	data      *reflect.Value
}

// NewComposite creates a new instance of the *Composite struct,
// validates and sets its Spec before returning it.
// Refer to SetSpec() for more information on Spec validation.
func NewComposite(spec *Spec) *Composite {
	f := &Composite{
		fieldsMap: map[int]struct{}{},
	}
	f.SetSpec(spec)
	return f
}

// Spec returns the receiver's spec.
func (f *Composite) Spec() *Spec {
	return f.spec
}

// SetSpec validates the spec and creates new instances of Fields defined
// in the specification.
// NOTE: Composite does not support padding. Therefore, users should
// only pass None or nil values for ths type. Passing any other value will
// result in a panic.
func (f *Composite) SetSpec(spec *Spec) {
	if err := validateCompositeSpec(spec); err != nil {
		panic(err)
	}
	f.spec = spec
	f.idToFieldMap = spec.CreateMessageFields()
	f.orderedSpecFieldIDs = orderedKeys(f.idToFieldMap)
}

// SetData traverses through fields provided in the data parameter matches them
// with their spec definition and calls SetData(...) on each spec field with the
// appropriate data field.
//
// A valid input is as follows:
//
//      type CompositeData struct {
//          F1 *String
//          F2 *String
//          F3 *Numeric
//          F4 *SubfieldCompositeData
//      }
//
func (f *Composite) SetData(data interface{}) error {
	dataStruct := reflect.ValueOf(data)
	if dataStruct.Kind() == reflect.Ptr || dataStruct.Kind() == reflect.Interface {
		// get the struct
		dataStruct = dataStruct.Elem()
	}

	if dataStruct.Kind() != reflect.Struct {
		return fmt.Errorf("failed to set data as struct is expected, got: %s", dataStruct.Kind())
	}

	f.data = &dataStruct
	return nil
}

// Pack deserialises data held by the receiver (via SetData)
// into bytes and returns an error on failure.
func (f *Composite) Pack() ([]byte, error) {
	packed, err := f.pack()
	if err != nil {
		return nil, err
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	return append(packedLength, packed...), nil
}

// Unpack takes in a byte array and serializes them into the receiver's
// subfields. An offset (unit depends on encoding and prefix values) is
// returned on success. A non-nil error is returned on failure.
func (f *Composite) Unpack(data []byte) (int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %v", err)
	}
	offset := f.spec.Pref.Length()

	// data is stripped of the prefix before it is provided to unpack().
	// Therefore, it is unaware of when to stop parsing unless we bound the
	// length of the slice by the data length.
	read, err := f.unpack(data[offset : offset+dataLen])
	if err != nil {
		return 0, err
	}
	if dataLen != read {
		return 0, fmt.Errorf("data length: %v does not match aggregate data read from decoded subfields: %v", dataLen, read)
	}

	return offset + read, nil
}

// SetBytes iterates over the receiver's subfields and unpacks them.
// Data passed into this method must consist of the necessary information to
// pack all subfields in full. However, unlike Unpack(), it requires the
// aggregate length of the subfields not to be encoded in the prefix.
func (f *Composite) SetBytes(data []byte) error {
	_, err := f.unpack(data)
	return err
}

// Bytes iterates over the receiver's subfields and packs them. The result
// does not incorporate the encoded aggregate length of the subfields in the
// prefix.
func (f *Composite) Bytes() ([]byte, error) {
	return f.pack()
}

// String iterates over the receiver's subfields, packs them and converts the
// result to a string. The result does not incorporate the encoded aggregate
// length of the subfields in the prefix.
func (f *Composite) String() (string, error) {
	b, err := f.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (f *Composite) MarshalJSON() ([]byte, error) {
	jsonData := OrderedMap(f.idToFieldMap)
	return json.Marshal(jsonData)
}

func (f *Composite) pack() ([]byte, error) {
	packed := []byte{}
	for _, id := range f.orderedSpecFieldIDs {
		specField := f.idToFieldMap[id]

		if f.data != nil {
			fieldName := fmt.Sprintf("F%d", id)
			// get the struct field
			dataField := f.data.FieldByName(fieldName)

			// no non-nil data field was found with this name
			if dataField == (reflect.Value{}) || dataField.IsNil() {
				continue
			}

			if err := specField.SetData(dataField.Interface()); err != nil {
				return nil, fmt.Errorf("failed to set data for field %d: %w", id, err)
			}
		}

		if f.spec.IDLength > 0 {
			idBytes, err := f.spec.Enc.Encode(idToBytes(f.spec.IDLength, id))
			if err != nil {
				return nil, fmt.Errorf("failed to convert subfield ID \"%s\" to int", idBytes)
			}
			packed = append(packed, idBytes...)
		}

		packedBytes, err := specField.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack subfield %d: %v", id, err)
		}
		packed = append(packed, packedBytes...)
	}
	return packed, nil
}

func (f *Composite) unpack(data []byte) (int, error) {
	if f.spec.IDLength > 0 {
		return f.unpackFieldsByID(data)
	}
	return f.unpackFields(data)
}

func (f *Composite) unpackFields(data []byte) (int, error) {
	offset := 0
	for _, id := range f.orderedSpecFieldIDs {
		specField := f.idToFieldMap[id]
		if err := f.setSubfieldData(id, specField); err != nil {
			return 0, err
		}
		read, err := specField.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %d: %v", id, err)
		}
		offset += read
	}
	return offset, nil
}

func (f *Composite) unpackFieldsByID(data []byte) (int, error) {
	offset := 0
	for offset < len(data) {
		idBytes, read, err := f.spec.Enc.Decode(data[offset:], f.spec.IDLength)
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield ID: %w", err)
		}

		id, err := strconv.Atoi(string(idBytes))
		if err != nil {
			return 0, fmt.Errorf("failed to convert subfield ID \"%s\" to int", string(idBytes))
		}

		specField, ok := f.idToFieldMap[id]
		if !ok {
			return 0, fmt.Errorf("failed to unpack subfield %d: field not defined in Spec", id)
		}
		offset += read

		if err := f.setSubfieldData(id, specField); err != nil {
			return 0, err
		}

		read, err = specField.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %d: %v", id, err)
		}
		offset += read
	}
	return offset, nil
}

func (f *Composite) setSubfieldData(id int, specField Field) error {
	if f.data == nil {
		return nil
	}

	fieldName := fmt.Sprintf("F%d", id)

	// get the struct field
	dataField := f.data.FieldByName(fieldName)

	// if data field was found with this name
	if dataField != (reflect.Value{}) {
		if dataField.IsNil() {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}
		if err := specField.SetData(dataField.Interface()); err != nil {
			return fmt.Errorf("failed to set data for field %d: %w", id, err)
		}
	}

	return nil
}

func validateCompositeSpec(spec *Spec) error {
	if spec.Pad != nil && spec.Pad != padding.None {
		return fmt.Errorf("Composite spec only supports nil or None padding values")
	}
	if spec.Enc == nil && spec.IDLength > 0 {
		return fmt.Errorf("Composite spec requires an Enc to be defined if IDLength > 0")
	}
	return nil
}

func orderedKeys(kvs map[int]Field) []int {
	keys := make([]int, 0)
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func idToBytes(length int, id int) []byte {
	idFmt := fmt.Sprintf("%%0%dd", length)
	return []byte(fmt.Sprintf(idFmt, id))
}
