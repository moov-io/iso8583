package field

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"io"
	"os"
	"reflect"
)

var _ Field = (*ConstructedTLV)(nil)
var _ json.Marshaler = (*ConstructedTLV)(nil)
var _ json.Unmarshaler = (*ConstructedTLV)(nil)

// A constructed BER-TLV data object consists of a tag, a length, and a value field
// composed of one or more BER-TLV data objects. A record in an AEF governed by
// this specification is a constructed BER-TLV data object.
//
//  Tag  Length  Primitive or constructed          Primitive or constructed
//  (T)   (L)       BER-TLV data object     ...      BER-TLV data object
//						number 1                          number n
//
// NOTE: Constructed TLV has not primitive value
//

type ConstructedTLV struct {
	spec *Spec

	orderedSpecFieldTags []string

	// stores all fields according to the spec
	subfields map[string]Field

	// tracks which subfields were set
	setSubfields map[string]struct{}
}

// NewConstructedTLV creates a new instance of the constructed tlv struct,
func NewConstructedTLV(spec *Spec) *ConstructedTLV {
	f := &ConstructedTLV{}
	f.SetSpec(spec)
	f.ConstructSubfields()

	return f
}

// NewConstructedTLVValue returns a instance of constructed tlv with value (raw bytes)
func NewConstructedTLVValue(val []byte) *ConstructedTLV {

	tlv := ConstructedTLV{}

	// Read subfields
	if err := tlv.SetBytes(val); err != nil {
		return nil
	}

	return &tlv
}

// NewConstructedTLVHexString returns a instance of constructed tlv with hex string
func NewConstructedTLVHexString(val string) *ConstructedTLV {

	value, err := encoding.BerTLVTag.Encode([]byte(val))
	if err != nil {
		return &ConstructedTLV{}
	}

	tlv := ConstructedTLV{}

	// Read subfields
	if err := tlv.SetBytes(value); err != nil {
		return nil
	}

	return &tlv
}

func (f *ConstructedTLV) ConstructSubfields() {
	if f.subfields == nil {
		f.subfields = CreateSubfields(f.spec)
	}
	f.setSubfields = make(map[string]struct{})
}

// Spec returns the receiver's spec.
func (f *ConstructedTLV) Spec() *Spec {
	return f.spec
}

// describe returns tlv tree with tag and value
func (f *ConstructedTLV) describe(output io.Writer, padLeft int) {

	if output == nil {
		output = os.Stdout
	}

	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			continue
		}

		if field.Spec() == nil || field.Spec().Tag == nil {
			continue
		}

		fmtStr := "%s\t:"
		for i := 0; i < padLeft; i++ {
			// spaces for tree levels
			fmtStr = "  " + fmtStr
		}

		if sub, ok := field.(*PrimitiveTLV); ok {
			raw, _ := sub.Bytes()
			fmt.Fprintf(output, fmtStr+" %X\n", field.Spec().Tag.Tag, raw)
		} else if sub, ok := field.(*ConstructedTLV); ok {
			fmt.Fprintf(output, fmtStr+"\n", field.Spec().Tag.Tag)
			sub.describe(output, padLeft+1)
		}
	}
}

// getSubfields returns the map of set sub fields
func (f *ConstructedTLV) getSubfields() map[string]Field {
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
func (f *ConstructedTLV) SetSpec(spec *Spec) {
	if err := validateConstructedTLVSpec(spec); err != nil {
		panic(err)
	}
	f.spec = spec
	f.orderedSpecFieldTags = orderedKeys(spec.Subfields, spec.Tag.Sort)
}

func (f *ConstructedTLV) Unmarshal(v interface{}) error {
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
func (f *ConstructedTLV) SetData(v interface{}) error {
	return f.Marshal(v)
}

func (f *ConstructedTLV) Marshal(v interface{}) error {
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
func (f *ConstructedTLV) Pack() ([]byte, error) {

	if f.spec.Tag == nil || f.spec.Tag.Enc == nil || f.spec.Pref == nil || f.spec.Enc == nil {
		return nil, fmt.Errorf("failed to pack tlv: invalid spec")
	}

	packed, err := f.pack()
	if err != nil {
		return nil, err
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	packed = append(packedLength, packed...)

	tagBytes := []byte(f.spec.Tag.Tag)
	if f.spec.Tag.Pad != nil {
		tagBytes = f.spec.Tag.Pad.Pad(tagBytes, f.spec.Tag.Length)
	}
	tagBytes, err = f.spec.Tag.Enc.Encode(tagBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode tlv Tag \"%v\"", tagBytes)
	}
	packed = append(tagBytes, packed...)

	return packed, nil
}

// Unpack takes in a byte array and serializes them into the receiver's
// subfields. An offset (unit depends on encoding and prefix values) is
// returned on success. A non-nil error is returned on failure.
func (f *ConstructedTLV) Unpack(data []byte) (int, error) {

	if f.spec.Tag == nil || f.spec.Tag.Enc == nil || f.spec.Pref == nil {
		return 0, fmt.Errorf("failed to unpack tlv: invalid spec")
	}

	offset := 0

	// 1. Read Tag

	tagBytes, read, err := f.spec.Tag.Enc.Decode(data[offset:], f.spec.Tag.Length)
	if err != nil {
		return 0, fmt.Errorf("failed to unpack subfield Tag: %w", err)
	}
	offset += read

	if f.spec.Tag.Pad != nil {
		tagBytes = f.spec.Tag.Pad.Unpad(tagBytes)
	}
	tag := string(tagBytes)

	if tag != f.spec.Tag.Tag {
		return 0, fmt.Errorf("tag mismatch: want to read %s, got %s", f.spec.Tag.Tag, tag)
	}

	// 2. Read Length

	dataLen, read, err := f.spec.Pref.DecodeLength(f.spec.Length, data[offset:])
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}
	offset += read

	// 3. Read Raw Value

	raw, read, err := f.spec.Enc.Decode(data[offset:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %w", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}
	offset += read

	// 4. Read subfields
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
func (f *ConstructedTLV) SetBytes(data []byte) error {
	_, err := f.unpack(data)
	return err
}

// Bytes iterates over the receiver's subfields and packs them. The result
// does not incorporate the encoded aggregate length of the subfields in the
// prefix.
func (f *ConstructedTLV) Bytes() ([]byte, error) {
	return f.pack()
}

// String iterates over the receiver's subfields, packs them and converts the
// result to a string. The result does not incorporate the encoded aggregate
// length of the subfields in the prefix.
func (f *ConstructedTLV) String() (string, error) {
	b, err := f.Bytes()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", b), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (f *ConstructedTLV) MarshalJSON() ([]byte, error) {
	jsonData := OrderedMap(f.getSubfields())
	return json.Marshal(jsonData)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// An error is thrown if the JSON consists of a subfield that has not
// been defined in the spec.
func (f *ConstructedTLV) UnmarshalJSON(b []byte) error {
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

func (f *ConstructedTLV) pack() ([]byte, error) {

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

func (f *ConstructedTLV) unpack(data []byte) (int, error) {

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

func validateConstructedTLVSpec(spec *Spec) error {
	if spec.Tag == nil || spec.Tag.Tag == "" {
		return fmt.Errorf("ConstructedTLV spec requires a Tag.tag value to be defined")
	}
	if spec.Pad != nil && spec.Pad != padding.None {
		return fmt.Errorf("ConstructedTLV spec only supports nil or None spec padding values")
	}
	if spec.Enc == nil {
		return fmt.Errorf("ConstructedTLV spec only supports a valid Enc value")
	}
	if spec.Tag != nil && spec.Tag.Enc == nil {
		return fmt.Errorf("ConstructedTLV spec requires a Tag.Enc")
	}
	if spec.Tag.Sort == nil {
		return fmt.Errorf("ConstructedTLV spec requires a Tag.Sort function to be defined")
	}
	return nil
}

// GetValue returns value of specified tag
func (f *ConstructedTLV) GetValue(tagHex string) ([]byte, error) {
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

// SetValue set value of tlv with specified tag
func (f *ConstructedTLV) SetValue(tagHex string, value []byte) error {
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
