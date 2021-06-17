package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
)

var _ Field = (*Subfields)(nil)

// Subfields is a wrapper object designed to hold ISO8583 subfields. It packs
// and unpacks subfields ordered by ID. The following is not supported by
// Subfields:
//  - Padding
//  - Encoding
//
// Responsibility for this is delegated recursively to the subfields
// themselves.
type Subfields struct {
	spec *Spec

	orderedFieldIDs []int
	idToFieldMap    map[int]Field

	data interface{}
}

// NewSubfields creates a new instance of the *Subfield struct,
// validates and sets its Spec before returning it.
// Refer to SetSpec() for more information on Spec validation.
func NewSubfields(spec *Spec) *Subfields {
	f := &Subfields{}
	f.SetSpec(spec)
	return f
}

// Spec returns the receiver's spec.
func (f *Subfields) Spec() *Spec {
	return f.spec
}

// SetSpec validates the spec and creates new instances of Fields defined
// in the specification.
// NOTE: Subfields do not support encoding and padding. Therefore, users should
// only pass None or nil values for these types. Passing any other value will
// result in a panic.
func (f *Subfields) SetSpec(spec *Spec) {
	if err := validateSubfieldsSpec(spec); err != nil {
		panic(err)
	}
	f.spec = spec
	f.idToFieldMap = spec.CreateMessageFields()
	f.orderedFieldIDs = orderedKeys(f.idToFieldMap)
}

// SetData traverses through fields provided in the data parameter matches them
// with their spec definition and calls SetData(...) on each spec field with the
// appropriate data field.
//
// A valid input is as follows:
//
//		type SubfieldData struct {
//			F1 *String
//			F2 *String
//			F3 *Numeric
//		}
//
func (f *Subfields) SetData(data interface{}) error {
	f.data = data

	if f.data == nil {
		return nil
	}

	dataStruct := reflect.ValueOf(data)
	if dataStruct.Kind() == reflect.Ptr || dataStruct.Kind() == reflect.Interface {
		// get the struct
		dataStruct = dataStruct.Elem()
	}

	if dataStruct.Kind() != reflect.Struct {
		return fmt.Errorf("failed to set data as struct is expected, got: %s", reflect.TypeOf(dataStruct).Kind())
	}

	for i, specField := range f.idToFieldMap {
		fieldName := fmt.Sprintf("F%d", i)

		// get the struct field
		dataField := dataStruct.FieldByName(fieldName)

		// no data field was found with this name
		if dataField == (reflect.Value{}) {
			continue
		}

		if dataField.IsNil() {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}
		if err := specField.SetData(dataField.Interface()); err != nil {
			return fmt.Errorf("failed to set data for field %d: %w", i, err)
		}
	}
	return nil
}

// Pack deserialises data held by the receiver (via SetData)
// into bytes and returns an error on failure.
func (f *Subfields) Pack() ([]byte, error) {
	packed := []byte{}
	for _, id := range f.orderedFieldIDs {
		packedBytes, err := f.idToFieldMap[id].Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack subfield %d: %v", id, err)
		}
		packed = append(packed, packedBytes...)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	return append(packedLength, packed...), nil
}

// Unpack takes in a byte array and serialises them into the receiver's
// subfields. An offset (unit depends on encoding and prefix values) is
// returned on success. A non-nil error is returned on failure.
func (f *Subfields) Unpack(data []byte) (int, error) {
	dataLen, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %v", err)
	}

	offset := f.spec.Pref.Length()
	read, err := f.unpack(data[offset:])
	if err != nil {
		return 0, err
	}
	if dataLen != read {
		return 0, fmt.Errorf("data length: %v does not match aggregate data read from decoded subfields: %v", dataLen, offset)
	}

	return read, nil
}

// SetBytes iterates over the receiver's subelements and unpacks them.
// Data passed into this method must consist of the necessary information to
// pack all subfields in full. However, unlike Unpack(), it requires the
// aggregate length of the subfields not to be encoded in the prefix.
func (f *Subfields) SetBytes(data []byte) error {
	_, err := f.unpack(data)
	return err
}

// Bytes iterates over the receiver's subelements and packs them. The result
// does not incorporate the encoded aggregate length of the subfields in the
// prefix.
func (f *Subfields) Bytes() ([]byte, error) {
	packed := []byte{}
	for _, id := range f.orderedFieldIDs {
		packedBytes, err := f.idToFieldMap[id].Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack struct subfield: %v", err)
		}
		packed = append(packed, packedBytes...)
	}
	return packed, nil
}

// String iterates over the receiver's subelements, packs them and converts the
// result to a string. The result does not incorporate the encoded aggregate
// length of the subfields in the prefix.
func (f *Subfields) String() (string, error) {
	b, err := f.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (f *Subfields) MarshalJSON() ([]byte, error) {
	jsonData := OrderedMap(f.idToFieldMap)
	return json.Marshal(jsonData)
}

func (f *Subfields) unpack(data []byte) (int, error) {
	offset := 0
	for _, id := range f.orderedFieldIDs {
		fl := f.idToFieldMap[id]
		read, err := fl.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %d: %v", id, err)
		}
		offset += read
	}
	return offset, nil
}

func validateSubfieldsSpec(spec *Spec) error {
	if spec.Enc != nil && spec.Enc != encoding.None {
		return fmt.Errorf("subfields spec only supports nil or None encoding values")
	}
	if spec.Pad != nil && spec.Pad != padding.None {
		return fmt.Errorf("subfields spec only supports nil or None padding values")
	}
	return nil
}

func orderedKeys(kvs map[int]Field) []int {
	keys := make([]int, 0)
	for k, _ := range kvs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
