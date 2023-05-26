package field

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"

	"github.com/moov-io/iso8583/utils"
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
// Because composite subfields may be encoded with different encodings, the
// Length field on the field.Spec struct is in bytes.
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
	spec   *Spec
	bitmap *Bitmap

	orderedSpecFieldTags []string

	// stores all fields according to the spec
	subfields map[string]Field

	// tracks which subfields were set
	setSubfields map[string]struct{}
}

// NewComposite creates a new instance of the *Composite struct,
// validates and sets its Spec before returning it.
// Refer to SetSpec() for more information on Spec validation.
func NewComposite(spec *Spec) *Composite {
	f := &Composite{}
	f.SetSpec(spec)
	f.ConstructSubfields()

	return f
}

// CompositeWithSubfields is used when composite field is created without
// calling NewComposite e.g. in iso8583.NewMessage(...)
type CompositeWithSubfields interface {
	ConstructSubfields()
}

func (f *Composite) ConstructSubfields() {
	if f.subfields == nil {
		f.subfields = CreateSubfields(f.spec)
	}
	f.setSubfields = make(map[string]struct{})
}

// Spec returns the receiver's spec.
func (f *Composite) Spec() *Spec {
	return f.spec
}

// getSubfields returns the map of set sub fields
func (f *Composite) getSubfields() map[string]Field {
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
func (f *Composite) SetSpec(spec *Spec) {
	if err := validateCompositeSpec(spec); err != nil {
		panic(err)
	}
	f.spec = spec

	var sortFn sort.StringSlice

	// When bitmap is defined, always order tags by int.
	if spec.Bitmap != nil {
		sortFn = sort.StringsByInt
	} else {
		sortFn = spec.Tag.Sort
	}

	f.orderedSpecFieldTags = orderedKeys(spec.Subfields, sortFn)
}

func (f *Composite) Unmarshal(v interface{}) error {
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
func (f *Composite) SetData(v interface{}) error {
	return f.Marshal(v)
}

// Marshal traverses through fields provided in the data parameter matches them
// with their spec definition and calls Marshal(...) on each spec field with the
// appropriate data field.
//
// A valid input is as follows:
//
//	type CompositeData struct {
//	    F1 *String
//	    F2 *String
//	    F3 *Numeric
//	    F4 *SubfieldCompositeData
//	}
func (f *Composite) Marshal(v interface{}) error {
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
func (f *Composite) Pack() ([]byte, error) {
	packed, err := f.pack()
	if err != nil {
		return nil, err
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

	isVariableLength := false
	if offset != 0 {
		isVariableLength = true
	}

	if offset+dataLen > len(data) {
		return 0, fmt.Errorf("not enough data to unpack, expected: %d, got: %d", offset+dataLen, len(data))
	}
	// data is stripped of the prefix before it is provided to unpack().
	// Therefore, it is unaware of when to stop parsing unless we bound the
	// length of the slice by the data length.
	read, err := f.unpack(data[offset:offset+dataLen], isVariableLength)
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

// Bitmap returns the parsed bitmap instantiated on the key "0" of the spec.
// In case the bitmap is not instantiated on the spec, returns nil.
func (f *Composite) Bitmap() *Bitmap {
	if f.bitmap != nil {
		return f.bitmap
	}

	if f.spec.Bitmap == nil {
		return nil
	}

	bitmap, ok := CreateSubfield(f.spec.Bitmap).(*Bitmap)
	if !ok {
		return nil
	}

	f.bitmap = bitmap

	return f.bitmap
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
	jsonData := OrderedMap(f.getSubfields())
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal map to bytes")
	}
	return bytes, nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// An error is thrown if the JSON consists of a subfield that has not
// been defined in the spec.
func (f *Composite) UnmarshalJSON(b []byte) error {
	var data map[string]json.RawMessage
	err := json.Unmarshal(b, &data)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to map")
	}

	for tag, rawMsg := range data {
		if _, ok := f.spec.Subfields[tag]; !ok && !f.skipUnknownTLVTags() {
			return fmt.Errorf("failed to unmarshal subfield %v: received subfield not defined in spec", tag)
		}

		subfield, ok := f.subfields[tag]
		if !ok {
			continue
		}

		if err := json.Unmarshal(rawMsg, subfield); err != nil {
			return utils.NewSafeErrorf(err, "failed to unmarshal subfield %v", tag)
		}

		f.setSubfields[tag] = struct{}{}
	}

	return nil
}

func (f *Composite) pack() ([]byte, error) {
	if f.Bitmap() != nil {
		return f.packByBitmap()
	}

	return f.packByTag()
}

func (f *Composite) packByBitmap() ([]byte, error) {
	f.Bitmap().Reset()

	// Set bitmap bits for all fields that are present.
	for id := range f.setSubfields {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("failed to pack composite: %w", err)
		}

		f.Bitmap().Set(idInt)
	}

	// pack bitmap.
	packed, err := f.Bitmap().Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack bitmap: %w", err)
	}

	// pack fields
	for _, i := range f.orderedSpecFieldTags {
		// If this ordered field is not set, continue to the next field.
		if _, ok := f.setSubfields[i]; !ok {
			continue
		}

		field, ok := f.subfields[i]
		if !ok {
			return nil, fmt.Errorf("failed to pack subfield %s: no specification found", i)
		}
		packedField, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack subfield %s (%s): %w", i, field.Spec().Description, err)
		}
		packed = append(packed, packedField...)
	}

	return packed, nil
}

func (f *Composite) packByTag() ([]byte, error) {
	packed := []byte{}
	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			return nil, fmt.Errorf("no subfield for tag %s", tag)
		}

		if _, set := f.setSubfields[tag]; !set {
			continue
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

func (f *Composite) unpack(data []byte, isVariableLength bool) (int, error) {
	if f.Bitmap() != nil {
		return f.unpackSubfieldsByBitmap(data)
	}
	if f.spec.Tag.Enc != nil {
		return f.unpackSubfieldsByTag(data)
	}
	return f.unpackSubfields(data, isVariableLength)
}

func (f *Composite) unpackSubfields(data []byte, isVariableLength bool) (int, error) {
	offset := 0
	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			continue
		}

		read, err := field.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}

		f.setSubfields[tag] = struct{}{}

		offset += read

		if isVariableLength && offset >= len(data) {
			break
		}
	}

	return offset, nil
}

func (f *Composite) unpackSubfieldsByBitmap(data []byte) (int, error) {
	var off int

	// Reset fields that were set.
	f.setSubfields = make(map[string]struct{})

	f.Bitmap().Reset()

	read, err := f.Bitmap().Unpack(data[off:])
	if err != nil {
		return 0, fmt.Errorf("failed to unpack bitmap: %w", err)
	}

	off += read

	for i := 1; i <= f.Bitmap().Len(); i++ {
		if f.Bitmap().IsSet(i) {
			iStr := strconv.Itoa(i)

			fl, ok := f.subfields[iStr]
			if !ok {
				return 0, fmt.Errorf("failed to unpack subfield %s: no specification found", iStr)
			}

			read, err = fl.Unpack(data[off:])
			if err != nil {
				return 0, fmt.Errorf("failed to unpack subfield %s (%s): %w", iStr, fl.Spec().Description, err)
			}

			f.setSubfields[iStr] = struct{}{}

			off += read
		}
	}

	return off, nil
}

// ignoredMaxLen is a constant meant to be used in encoders that don't use the maxLength argument during
// length decoding.
const ignoredMaxLen int = 0

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
			if f.skipUnknownTLVTags() {
				// Obtain the length of the unknown tag and add it to the offset.
				// Because BER-TLV lengths are decoded dynamically, the maxLen method argument is ignored.
				fieldLength, readed, err := prefix.BerTLV.DecodeLength(ignoredMaxLen, data[offset:])
				if err != nil {
					return 0, err
				}
				offset += fieldLength + readed
				continue
			}

			return 0, fmt.Errorf("failed to unpack subfield %v: field not defined in Spec", tag)
		}

		field, ok := f.subfields[tag]
		if !ok {
			continue
		}

		read, err = field.Unpack(data[offset:])
		if err != nil {
			return 0, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}

		f.setSubfields[tag] = struct{}{}

		offset += read
	}
	return offset, nil
}

func (f *Composite) skipUnknownTLVTags() bool {
	return f.spec.Tag != nil && f.spec.Tag.SkipUnknownTLVTags && f.spec.Tag.Enc == encoding.BerTLVTag
}

func validateCompositeSpec(spec *Spec) error {
	if spec.Enc != nil {
		return fmt.Errorf("Composite spec only supports a nil Enc value")
	}
	if spec.Pad != nil && spec.Pad != padding.None {
		return fmt.Errorf("Composite spec only supports nil or None spec padding values")
	}
	if (spec.Bitmap == nil && spec.Tag == nil) || (spec.Bitmap != nil && spec.Tag != nil) {
		return fmt.Errorf("Composite spec only supports a definition of Bitmap or Tag, can't stand both or neither")
	}

	// If bitmap is defined, validates subfields keys.
	// spec.Tag is not validated.
	if spec.Bitmap != nil {
		if !spec.Bitmap.spec.DisableAutoExpand {
			return fmt.Errorf("Composite spec only supports a bitmap with 'DisableAutoExpand = true'")
		}

		for key := range spec.Subfields {
			parsedKey, err := strconv.Atoi(key)
			if err != nil {
				return fmt.Errorf("error parsing key from bitmapped subfield definition: %w", err)
			}

			if parsedKey <= 0 {
				return fmt.Errorf("Composite spec only supports integers greater than 0 as keys for bitmapped subfields definition")
			}
		}

		return nil
	}

	// Validate spec.Tag.
	if spec.Tag.Sort == nil {
		return fmt.Errorf("Composite spec requires a Tag.Sort function to define a Tag")
	}
	if spec.Tag.Enc == nil && spec.Tag.Length > 0 {
		return fmt.Errorf("Composite spec requires a Tag.Enc to be defined if Tag.Length > 0")
	}

	return nil
}

func orderedKeys(kvs map[string]Field, sorter sort.StringSlice) []string {
	keys := make([]string, 0, len(kvs))
	for k := range kvs {
		keys = append(keys, k)
	}
	sorter(keys)
	return keys
}

var fieldNameTagRe = regexp.MustCompile(`^F.+$`)

// getFieldIndexOrTag returns index or tag of the field. First, it checks the
// field name. If it does not match F.+ pattern, it checks value of `index`
// tag.  If empty string, then index/tag was not found for the field.
func getFieldIndexOrTag(field reflect.StructField) (string, error) {
	dataFieldName := field.Name

	if fieldIndex := field.Tag.Get("index"); fieldIndex != "" {
		return fieldIndex, nil
	}

	if len(dataFieldName) > 0 && fieldNameTagRe.MatchString(dataFieldName) {
		return dataFieldName[1:], nil
	}

	return "", nil
}
