package field

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/padding"
)

var _ Field = (*EMV)(nil)
var _ json.Marshaler = (*EMV)(nil)
var _ json.Unmarshaler = (*EMV)(nil)

// EMV is a payment method based on a technical standard for smart payment cards and for payment terminals
// and automated teller machines which can accept them. EMV stands for "Europay, Mastercard, and Visa",
// the three companies that created the standard.
//
// EMV cards are smart cards, also called chip cards, integrated circuit cards, or IC cards,
// which store their data on integrated circuit chips, in addition to magnetic stripes for backward compatibility.
//
// ISO 8583 – Field or DE 55 has the EMV data encoded in the authorization message.
//

type EMV struct {
	spec *Spec

	orderedSpecFieldTags []string

	// stores all fields according to the spec
	subfields map[string]Field

	// tracks which subfields were set
	setSubfields map[string]struct{}
}

// NewEMV creates a new instance of the constructed tlv struct,
func NewEMV(spec *Spec) *EMV {
	f := &EMV{}
	f.SetSpec(spec)
	f.ConstructSubfields()

	return f
}

func (f *EMV) ConstructSubfields() {
	if f.subfields == nil {
		f.subfields = CreateSubfields(f.spec)
	}
	f.setSubfields = make(map[string]struct{})
}

// Spec returns the receiver's spec.
func (f *EMV) Spec() *Spec {
	return f.spec
}

// getSubfields returns the map of set sub fields
func (f *EMV) getSubfields() map[string]Field {
	fields := map[string]Field{}
	for i := range f.setSubfields {
		fields[i] = f.subfields[i]
	}
	return fields
}

// SetSpec validates the spec and creates new instances of Subfields defined
// in the specification.
// NOTE: Composite does not support padding on the base spec. Therefore, users
// should only pass None or nil values for ths type. Passing any other value
// will result in a panic.
func (f *EMV) SetSpec(spec *Spec) {
	if err := validateEMVSpec(spec); err != nil {
		panic(err)
	}
	f.spec = spec
	f.orderedSpecFieldTags = orderedKeys(spec.Subfields, spec.Tag.Sort)
}

func (f *EMV) Unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("data is not a pointer or nil")
	}

	// get the struct from the pointer
	dataStruct := rv.Elem()

	if dataStruct.Kind() != reflect.Struct {
		return errors.New("data is not a struct")
	}

	// iterate over struct fields
	for i := 0; i < dataStruct.NumField(); i++ {
		indexOrTag, err := getFieldIndexOrTag(dataStruct.Type().Field(i))
		if err != nil {
			return fmt.Errorf("getting field %d index: %w", i, err)
		}

		// skip field without index
		if indexOrTag == "" {
			continue
		}

		messageField, ok := f.subfields[indexOrTag]
		if !ok {
			continue
		}

		// unmarshal only subfield that has the value set
		if _, set := f.setSubfields[indexOrTag]; !set {
			continue
		}

		dataField := dataStruct.Field(i)
		if dataField.IsNil() {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}

		err = messageField.Unmarshal(dataField.Interface())
		if err != nil {
			return fmt.Errorf("failed to get data from field %s: %w", indexOrTag, err)
		}
	}

	return nil
}

// Deprecated. Use Marshal instead
func (f *EMV) SetData(v interface{}) error {
	return f.Marshal(v)
}

func (f *EMV) Marshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("data is not a pointer or nil")
	}

	// get the struct from the pointer
	dataStruct := rv.Elem()

	if dataStruct.Kind() != reflect.Struct {
		return errors.New("data is not a struct")
	}

	// iterate over struct fields
	for i := 0; i < dataStruct.NumField(); i++ {
		indexOrTag, err := getFieldIndexOrTag(dataStruct.Type().Field(i))
		if err != nil {
			return fmt.Errorf("getting field %d index: %w", i, err)
		}

		// skip field without index
		if indexOrTag == "" {
			continue
		}

		messageField, ok := f.subfields[indexOrTag]
		if !ok {
			continue
		}

		dataField := dataStruct.Field(i)
		if dataField.IsNil() {
			continue
		}

		err = messageField.Marshal(dataField.Interface())
		if err != nil {
			return fmt.Errorf("failed to set data from field %s: %w", indexOrTag, err)
		}

		f.setSubfields[indexOrTag] = struct{}{}
	}

	return nil
}

// Pack deserialises data held by the receiver (via SetData)
// into bytes and returns an error on failure.
func (f *EMV) Pack() ([]byte, error) {

	packed, err := f.pack()
	if err != nil {
		return nil, err
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	packed = append(packedLength, packed...)

	return packed, nil
}

// Unpack takes in a byte array and serializes them into the receiver's
// subfields. An offset (unit depends on encoding and prefix values) is
// returned on success. A non-nil error is returned on failure.
func (f *EMV) Unpack(data []byte) (int, error) {

	offset := 0

	dataLen, read, err := f.spec.Pref.DecodeLength(f.spec.Length, data[offset:])
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}
	offset += read

	raw, read, err := f.spec.Enc.Decode(data[offset:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %w", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}
	offset += read

	_, err = f.unpack(raw)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

// SetBytes iterates over the receiver's subfields and unpacks them.
// Data passed into this method must consist of the necessary information to
// pack all subfields in full. However, unlike Unpack(), it requires the
// aggregate length of the subfields not to be encoded in the prefix.
func (f *EMV) SetBytes(data []byte) error {
	_, err := f.unpack(data)
	return err
}

// Bytes iterates over the receiver's subfields and packs them. The result
// does not incorporate the encoded aggregate length of the subfields in the
// prefix.
func (f *EMV) Bytes() ([]byte, error) {
	return f.pack()
}

// String iterates over the receiver's subfields, packs them and converts the
// result to a string. The result does not incorporate the encoded aggregate
// length of the subfields in the prefix.
func (f *EMV) String() (string, error) {
	b, err := f.Bytes()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", b), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (f *EMV) MarshalJSON() ([]byte, error) {
	jsonData := OrderedMap(f.getSubfields())
	return json.Marshal(jsonData)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// An error is thrown if the JSON consists of a subfield that has not
// been defined in the spec.
func (f *EMV) UnmarshalJSON(b []byte) error {
	var data map[string]json.RawMessage
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	for tag, rawMsg := range data {
		if _, ok := f.spec.Subfields[tag]; !ok {
			return fmt.Errorf("failed to unmarshal subfield %v: received subfield not defined in spec", tag)
		}

		subfield, ok := f.subfields[tag]
		if !ok {
			continue
		}

		if err := json.Unmarshal(rawMsg, subfield); err != nil {
			return fmt.Errorf("failed to unmarshal subfield %v: %w", tag, err)
		}

		f.setSubfields[tag] = struct{}{}
	}

	return nil
}

func (f *EMV) pack() ([]byte, error) {

	packed := []byte{}

	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			return nil, fmt.Errorf("no subfield for tag %s", tag)
		}

		if _, set := f.setSubfields[tag]; !set {
			continue
		}

		packedBytes, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack subfield %v: %w", tag, err)
		}
		packed = append(packed, packedBytes...)

	}

	return packed, nil
}

func (f *EMV) unpack(data []byte) (int, error) {

	offset := 0

	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			return 0, fmt.Errorf("no subfield for tag %s", tag)
		}

		if field.Spec() == nil || field.Spec().Tag == nil {
			return 0, fmt.Errorf("no spec for subfield")
		}

		read, err := field.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %v: %w", field.Spec().Tag.Tag, err)
		}

		f.setSubfields[tag] = struct{}{}

		offset += read
	}

	return offset, nil
}

func validateEMVSpec(spec *Spec) error {
	if spec.Pad != nil && spec.Pad != padding.None {
		return fmt.Errorf("EMV spec only supports nil or None spec padding values")
	}
	if spec.Enc == nil {
		return fmt.Errorf("EMV spec only supports a valid Enc value")
	}
	if spec.Pref == nil {
		return fmt.Errorf("EMV spec only supports a valid pref")
	}
	if spec.Tag.Sort == nil {
		return fmt.Errorf("EMV spec requires a Tag.Sort function to be defined")
	}
	return nil
}

// GetValue returns value of specified tag
func (f *EMV) GetValue(tagHex string) ([]byte, error) {
	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			return nil, fmt.Errorf("unabled to find the tag %s", tagHex)
		}

		if sub, ok := field.(*ConstructedTLV); ok {
			value, getErr := sub.GetValue(tagHex)
			if getErr == nil {
				return value, nil
			}
		}

		if field.Spec() == nil || field.Spec().Tag == nil || field.Spec().Tag.Tag == "" {
			return nil, fmt.Errorf("no spec for the tag  %s", tagHex)
		}

		if field.Spec().Tag.Tag == tagHex {
			return field.Bytes()
		}
	}

	return nil, fmt.Errorf("unabled to find the tag %s", tagHex)
}

// SetValue set value of sub tlv with specified tag
func (f *EMV) SetValue(tagHex string, value []byte) error {
	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			return fmt.Errorf("unabled to find the tag %s", tagHex)
		}

		if sub, ok := field.(*ConstructedTLV); ok {
			getErr := sub.SetValue(tagHex, value)
			if getErr == nil {
				return nil
			}
		}

		if field.Spec() == nil || field.Spec().Tag == nil || field.Spec().Tag.Tag == "" {
			return fmt.Errorf("no spec for the tag  %s", tagHex)
		}

		if field.Spec().Tag.Tag == tagHex {
			return field.SetBytes(value)
		}
	}

	return fmt.Errorf("unabled to find the tag %s", tagHex)
}
