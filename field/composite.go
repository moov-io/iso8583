package field

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/sort"
)

var _ Field = (*Composite)(nil)
var _ json.Marshaler = (*Composite)(nil)
var _ json.Unmarshaler = (*Composite)(nil)

// Composite is a wrapper object designed to hold ISO8583 TLVs, subfields and
// subelements. Because Composite handles both of these usecases generically,
// we refer to them collectively as 'subfields' throughout the receiver's
// documentation and error messages. These subfields are defined using the
// 'Subfields' field on the field.Spec struct.
//
// Composite handles aggregate fields of the following format:
// - Length (if variable)
// - []Subfield
//
// Where the subfield structure is assumed to be as follows:
// - Subfield Tag (if Composite.Spec().Tag.Enc != nil)
// - Subfield Length (if variable)
// - Subfield Data (or Value in the case of TLVs)
//
// Composite behaves in two modes depending on whether subfield Tags need to be
// explicitly handled or not. This is configured by setting Spec.Tag.Length.
//
// When populated, Composite handles the packing and unpacking subfield Tags on
// their behalf. However, responsibility for packing and unpacking both the
// length and data of subfields is delegated to the subfields themselves.
// Therefore, their specs should be configured to handle such behavior.
//
// If Spec.Tag.Length != nil, Composite leverages Spec.Tag.Enc to unpack subfields
// regardless of order based on their Tags. Similarly, it is also able to handle
// non-present subfields by relying on the existence of their Tags.
//
// If Spec.Tag.Length == nil, Composite only unpacks subfields ordered by Tag.
// This order is determined by calling Spec.Tag.Sort on a slice of all subfield
// keys defined in the spec. The absence of Tags in the payload means that the
// receiver is not able to handle non-present subfields either.
//
// Tag.Pad should be used to set the padding direction and type of the Tag in
// situations when tags hold leading characters e.g. '003'. Both the data struct
// and the Spec.Subfields map may then omit those padded characters in their
// respective definitions.
//
// For the sake of determinism, packing of subfields is executed in order of Tag
// (using Spec.Tag.Sort) regardless of the value of Spec.Tag.Length.
type Composite struct {
	spec *Spec

	orderedSpecFieldTags []string
	tagToSubfieldMap     map[string]Field

	data *reflect.Value
}

// NewComposite creates a new instance of the *Composite struct,
// validates and sets its Spec before returning it.
// Refer to SetSpec() for more information on Spec validation.
func NewComposite(spec *Spec) *Composite {
	f := &Composite{}
	f.SetSpec(spec)
	return f
}

// Spec returns the receiver's spec.
func (f *Composite) Spec() *Spec {
	return f.spec
}

// SetSpec validates the spec and creates new instances of Subfields defined
// in the specification.
// NOTE: Composite does not support padding on the base spec. Therefore, users
// should only pass None or nil values for ths type. Passing any other value
// will result in a panic.
func (f *Composite) SetSpec(spec *Spec) {
	if err := validateCompositeSpec(spec); err != nil {
		panic(err)
	}
	f.spec = spec
	f.tagToSubfieldMap = map[string]Field{}
	f.orderedSpecFieldTags = orderedKeys(spec.Subfields, spec.Tag.Sort)
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

	if len(packed) == 0 {
		return []byte{}, nil
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

// Unpack takes in a byte array and serializes them into the receiver's
// subfields. An offset (unit depends on encoding and prefix values) is
// returned on success. A non-nil error is returned on failure.
func (f *Composite) Unpack(data []byte) (int, error) {
	dataLen, offset, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}

	hasPrefix := false
	if offset != 0 {
		hasPrefix = true
	}

	// data is stripped of the prefix before it is provided to unpack().
	// Therefore, it is unaware of when to stop parsing unless we bound the
	// length of the slice by the data length.
	read, err := f.unpack(data[offset:offset+dataLen], hasPrefix)
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
	_, err := f.unpack(data, false)
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
	// We pack the field to populate f.tagToSubfieldMap
	if _, err := f.Pack(); err != nil {
		return nil, err
	}
	jsonData := OrderedMap(f.tagToSubfieldMap)
	return json.Marshal(jsonData)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// An error is thrown if the JSON consists of a subfield that has not
// been defined in the spec.
func (f *Composite) UnmarshalJSON(b []byte) error {
	var data map[string]json.RawMessage
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	for tag, rawMsg := range data {
		if _, ok := f.spec.Subfields[tag]; !ok {
			return fmt.Errorf("failed to unmarshal subfield %v: received subfield not defined in spec", tag)
		}
		subfield := f.createAndSetSubfield(tag)

		if err = f.setSubfieldData(tag, subfield); err != nil {
			return err
		}
		if err := json.Unmarshal(rawMsg, subfield); err != nil {
			return fmt.Errorf("failed to unmarshal subfield %v: %w", tag, err)
		}
	}

	return nil
}

func (f *Composite) pack() ([]byte, error) {
	packed := []byte{}
	for _, tag := range f.orderedSpecFieldTags {
		var field Field
		if f.data != nil {
			fieldName := fmt.Sprintf("F%v", tag)
			// get the struct field
			dataField := f.data.FieldByName(fieldName)

			// no non-nil data field was found with this name
			if dataField == (reflect.Value{}) || dataField.IsNil() {
				continue
			}
			field = f.createAndSetSubfield(tag)

			if err := field.SetData(dataField.Interface()); err != nil {
				return nil, fmt.Errorf("failed to set data for field %v: %w", tag, err)
			}
		} else {
			var ok bool
			field, ok = f.tagToSubfieldMap[tag]
			if !ok {
				continue
			}
		}

		if f.spec.Tag != nil && f.spec.Tag.Enc != nil {
			tagBytes := []byte(tag)
			if f.spec.Tag.Pad != nil {
				tagBytes = f.spec.Tag.Pad.Pad(tagBytes, f.spec.Tag.Length)
			}
			tagBytes, err := f.spec.Tag.Enc.Encode(tagBytes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert subfield Tag \"%v\" to int", tagBytes)
			}
			packed = append(packed, tagBytes...)
		}

		packedBytes, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack subfield %v: %w", tag, err)
		}
		packed = append(packed, packedBytes...)
	}
	return packed, nil
}

func (f *Composite) unpack(data []byte, hasPrefix bool) (int, error) {
	if f.spec.Tag.Enc != nil {
		return f.unpackSubfieldsByTag(data)
	}
	return f.unpackSubfields(data, hasPrefix)
}

func (f *Composite) unpackSubfields(data []byte, hasPrefix bool) (int, error) {
	offset := 0
	for _, tag := range f.orderedSpecFieldTags {
		field := f.createAndSetSubfield(tag)
		if err := f.setSubfieldData(tag, field); err != nil {
			return 0, err
		}
		read, err := field.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}
		offset += read

		if hasPrefix && offset >= len(data) {
			break
		}
	}
	return offset, nil
}

func (f *Composite) unpackSubfieldsByTag(data []byte) (int, error) {
	offset := 0
	for offset < len(data) {
		tagBytes, read, err := f.spec.Tag.Enc.Decode(data[offset:], f.spec.Tag.Length)
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield Tag: %w", err)
		}
		offset += read

		if f.spec.Tag.Pad != nil {
			tagBytes = f.spec.Tag.Pad.Unpad(tagBytes)
		}
		tag := string(tagBytes)
		if _, ok := f.spec.Subfields[tag]; !ok {
			return 0, fmt.Errorf("failed to unpack subfield %v: field not defined in Spec", tag)
		}
		field := f.createAndSetSubfield(tag)

		if err := f.setSubfieldData(tag, field); err != nil {
			return 0, err
		}

		read, err = field.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}
		offset += read
	}
	return offset, nil
}

func (f *Composite) setSubfieldData(tag string, specField Field) error {
	if f.data == nil {
		return nil
	}

	fieldName := fmt.Sprintf("F%v", tag)

	// get the struct field
	dataField := f.data.FieldByName(fieldName)

	// if data field was found with this name
	if dataField != (reflect.Value{}) {
		if dataField.IsNil() {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}
		if err := specField.SetData(dataField.Interface()); err != nil {
			return fmt.Errorf("failed to set data for field %v: %w", tag, err)
		}
	}

	return nil
}

func (f *Composite) createAndSetSubfield(tag string) Field {
	field := CreateSubfield(f.spec.Subfields[tag])
	f.tagToSubfieldMap[tag] = field
	return field
}

func validateCompositeSpec(spec *Spec) error {
	if spec.Tag == nil || spec.Tag.Sort == nil {
		return fmt.Errorf("Composite spec requires a Tag.Sort function to be defined")
	}
	if spec.Pad != nil && spec.Pad != padding.None {
		return fmt.Errorf("Composite spec only supports nil or None spec padding values")
	}
	if spec.Enc != nil {
		return fmt.Errorf("Composite spec only supports a nil Enc value")
	}
	if spec.Tag != nil && spec.Tag.Enc == nil && spec.Tag.Length > 0 {
		return fmt.Errorf("Composite spec requires a Tag.Enc to be defined if Tag.Length > 0")
	}
	if spec.Tag.Sort == nil {
		return fmt.Errorf("Composite spec requires a Tag.Sort function to be defined")
	}
	return nil
}

func orderedKeys(kvs map[string]Field, sorter sort.StringSlice) []string {
	keys := make([]string, 0)
	for k := range kvs {
		keys = append(keys, k)
	}
	sorter(keys)
	return keys
}
