package field

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/moov-io/iso8583/encoding"
	iso8583errors "github.com/moov-io/iso8583/errors"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/moov-io/iso8583/utils"
)

var (
	_ Field            = (*Composite)(nil)
	_ json.Marshaler   = (*Composite)(nil)
	_ json.Unmarshaler = (*Composite)(nil)
)

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
	spec         *Spec
	cachedBitmap *Bitmap

	orderedSpecFieldTags []string

	// mu is used to synchronize access to the subfields and
	// setSubfields maps when the composite is used concurrently
	mu sync.Mutex

	// stores all fields according to the spec
	subfields map[string]Field
}

// NewComposite creates a new instance of the *Composite struct,
// validates and sets its Spec before returning it.
// Refer to SetSpec() for more information on Spec validation.
func NewComposite(spec *Spec) *Composite {
	f := &Composite{
		subfields: make(map[string]Field),
	}
	f.SetSpec(spec)

	return f
}

func (c *Composite) NewInstance() Field {
	return &Composite{
		spec:      c.spec, // spec is validated already
		subfields: make(map[string]Field),
	}
}

// Spec returns the receiver's spec.
func (f *Composite) Spec() *Spec {
	return f.spec
}

// GetSubfields returns the map of set sub fields. The returned map is a copy
// of the internal map, but the fields themselves are live references.
func (f *Composite) GetSubfields() map[string]Field {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.getSubfields()
}

// getSubfields returns the map of set sub fields, it should be called
// only when the mutex is locked
func (f *Composite) getSubfields() map[string]Field {
	return maps.Clone(f.subfields)
}

// SetSpec validates the spec and creates new instances of Subfields defined
// in the specification.
// NOTE: Composite does not support padding on the base spec. Therefore, users
// should only pass None or nil values for this type. Passing any other value
// will result in a panic.
func (f *Composite) SetSpec(spec *Spec) {
	if err := spec.Validate(); err != nil {
		panic(err) //nolint:forbidigo,nolintlint // as specs mostly static, we panic on spec validation errors
	}
	f.spec = spec
}

func (f *Composite) Unmarshal(v any) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("data is not a pointer or nil")
	}

	// get the struct from the pointer
	dataStruct := rv.Elem()
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("data is not a struct")
	}

	// iterate over struct fields
	for i := range dataStruct.NumField() {
		indexTag := NewIndexTag(dataStruct.Type().Field(i))
		if indexTag.Tag == "" {
			continue
		}

		messageField, ok := f.subfields[indexTag.Tag]
		if !ok {
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
				return fmt.Errorf("unmarshalling field %s: %w", indexTag.Tag, err)
			}
		case reflect.Slice:
			// Pass reflect.Value for slices so they can be modified
			err := messageField.Unmarshal(dataField)
			if err != nil {
				return fmt.Errorf("unmarshalling field %s: %w", indexTag.Tag, err)
			}
		default: // Native types
			err := messageField.Unmarshal(dataField)
			if err != nil {
				return fmt.Errorf("unmarshalling field %s: %w", indexTag.Tag, err)
			}
		}
	}

	return nil
}

// Deprecated. Use Marshal instead
func (f *Composite) SetData(v any) error {
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
func (f *Composite) Marshal(v any) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return errors.New("data is not a pointer")
	}

	elemType := rv.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return errors.New("data must be a pointer to struct")
	}

	// If nil, create a new instance of the struct
	if rv.IsNil() {
		rv = reflect.New(elemType)
	}

	// get the struct from the pointer
	dataStruct := rv.Elem()

	// iterate over struct fields
	for i := range dataStruct.NumField() {
		indexTag := NewIndexTag(dataStruct.Type().Field(i))
		if indexTag.Tag == "" {
			continue
		}

		dataField := dataStruct.Field(i)
		if dataField.IsZero() && !indexTag.KeepZero {
			continue
		}

		messageField, err := f.getOrCreateField(indexTag.Tag)
		if err != nil { // at the moment, getOrCreateField only returns error if field not in spec
			continue
		}

		err = messageField.Marshal(dataField.Interface())
		if err != nil {
			return fmt.Errorf("marshalling field %s: %w", indexTag.Tag, err)
		}
	}

	return nil
}

func (f *Composite) getOrCreateField(id string) (Field, error) {
	field := f.subfields[id]
	if field != nil {
		return field, nil
	}

	field, err := f.createField(id)
	if err != nil {
		return nil, fmt.Errorf("creating field %s: %w", id, err)
	}

	return field, nil
}

func (f *Composite) createField(id string) (Field, error) {
	specField, ok := f.Spec().Subfields[id]
	if !ok {
		return nil, fmt.Errorf("field %s is not defined in the spec", id)
	}

	field := NewInstanceOf(specField)
	f.subfields[id] = field

	return field, nil
}

// Pack deserialises data held by the receiver (via SetData)
// into bytes and returns an error on failure.
func (f *Composite) Pack() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

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
	f.mu.Lock()
	defer f.mu.Unlock()

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
	read, err := f.wrapErrorUnpack(data[offset:offset+dataLen], isVariableLength)
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
	f.mu.Lock()
	defer f.mu.Unlock()

	_, err := f.wrapErrorUnpack(data, false)
	return err
}

// Bytes iterates over the receiver's subfields and packs them. The result
// does not incorporate the encoded aggregate length of the subfields in the
// prefix.
func (f *Composite) Bytes() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.pack()
}

// Bitmap returns the parsed bitmap instantiated on the key "0" of the spec.
// In case the bitmap is not instantiated on the spec, returns nil.
func (f *Composite) Bitmap() *Bitmap {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.bitmap()
}

func (f *Composite) bitmap() *Bitmap {
	if f.cachedBitmap != nil {
		return f.cachedBitmap
	}

	if f.spec.Bitmap == nil {
		return nil
	}

	// we already know that spec.Bitmap is of type *Bitmap
	// and it presents the spec
	//nolint:forcetypeassert
	bitmap := f.spec.Bitmap.NewInstance().(*Bitmap)
	bitmap.Reset()

	f.cachedBitmap = bitmap

	return f.cachedBitmap
}

func (f *Composite) isWithBitmap() bool {
	return f.spec.Bitmap != nil
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
	f.mu.Lock()
	defer f.mu.Unlock()

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
	f.mu.Lock()
	defer f.mu.Unlock()

	var data map[string]json.RawMessage
	err := json.Unmarshal(b, &data)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to map")
	}

	for tag, rawMsg := range data {
		subfield, err := f.getOrCreateField(tag)
		if err != nil { // at the moment, getOrCreateField only returns error if field not in spec
			if !f.skipUnknownTLVTags() {
				return fmt.Errorf("getting or creating subfield %s: %w", tag, err)
			}
		}

		// subfield is not defined in spec and we skip unknown TLV tags
		if subfield == nil {
			continue
		}

		if err := json.Unmarshal(rawMsg, subfield); err != nil {
			return utils.NewSafeErrorf(err, "failed to unmarshal subfield %v", tag)
		}
	}

	return nil
}

func (f *Composite) pack() ([]byte, error) {
	// we can validate specs here
	if f.isWithBitmap() {
		return f.packWithBitmap()
	}

	return f.packByTag()
}

func (f *Composite) packWithBitmap() ([]byte, error) {
	f.bitmap().Reset()

	var packedFields []byte

	for _, id := range orderedKeys(f.subfields, sort.StringsByInt) {
		// field must be there as we got it from the subfields map
		field := f.subfields[id]

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("converting id %s to int: %w", id, err)
		}

		// set bitmap bit for this field
		f.bitmap().Set(idInt)

		packedField, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack subfield %s (%s): %w", id, field.Spec().Description, err)
		}

		packedFields = append(packedFields, packedField...)
	}

	// pack bitmap.
	packedBitmap, err := f.bitmap().Pack()
	if err != nil {
		return nil, fmt.Errorf("packing bitmap: %w", err)
	}

	return append(packedBitmap, packedFields...), nil
}

func (f *Composite) packByTag() ([]byte, error) {
	packed := []byte{}

	if f.spec.Tag == nil {
		return nil, errors.New("cannot pack composite field by tag when Tag spec is not defined")
	}

	for _, tag := range orderedKeys(f.subfields, f.spec.Tag.Sort) {
		field := f.subfields[tag]

		if f.spec.Tag.Enc != nil {
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

// wrapErrorUnpack calls the core unpacking logic and wraps any
// errors in a *UnpackError. It assumes that the mutex is already
// locked by the caller.
func (f *Composite) wrapErrorUnpack(src []byte, isVariableLength bool) (int, error) {
	offset, tagID, err := f.unpack(src, isVariableLength)
	if err != nil {
		return offset, &iso8583errors.UnpackError{
			Err:        err,
			FieldID:    tagID,
			RawMessage: src,
		}
	}

	return offset, nil
}

func (f *Composite) unpack(data []byte, isVariableLength bool) (int, string, error) {
	if f.isWithBitmap() {
		n, s, err := f.unpackSubfieldsByBitmap(data)
		if err != nil {
			return n, s, fmt.Errorf("unpacking subfields by bitmap: %w", err)

		}

		return n, s, nil
	}

	if f.spec.Tag.Enc != nil {
		n, s, err := f.unpackSubfieldsByTag(data)
		if err != nil {
			return n, s, fmt.Errorf("unpacking subfields by tag: %w", err)
		}

		return n, s, nil
	}

	n, s, err := f.unpackSubfields(data, isVariableLength)
	if err != nil {
		return n, s, fmt.Errorf("unpacking subfields: %w", err)
	}

	return n, s, nil
}

func (f *Composite) unpackSubfields(data []byte, isVariableLength bool) (int, string, error) {
	// reset subfields
	f.subfields = make(map[string]Field)

	offset := 0

	for _, tag := range orderedKeys(f.spec.Subfields, f.spec.Tag.Sort) {
		field := NewInstanceOf(f.spec.Subfields[tag])
		f.subfields[tag] = field

		read, err := field.Unpack(data[offset:])
		if err != nil {
			return 0, tag, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}

		offset += read

		if isVariableLength && offset >= len(data) {
			break
		}
	}

	return offset, "", nil
}

func (f *Composite) unpackSubfieldsByBitmap(data []byte) (int, string, error) {
	// reset subfields and bitmap
	f.subfields = make(map[string]Field)
	f.bitmap().Reset()

	offset := 0

	read, err := f.bitmap().Unpack(data[offset:])
	if err != nil {
		return 0, "", fmt.Errorf("failed to unpack bitmap: %w", err)
	}

	offset += read

	for i := 1; i <= f.bitmap().Len(); i++ {
		if f.bitmap().IsSet(i) {
			idx := strconv.Itoa(i)

			fl, err := f.createField(idx)
			if err != nil {
				return 0, idx, fmt.Errorf("getting or creating subfield %s: %w", idx, err)
			}

			read, err = fl.Unpack(data[offset:])
			if err != nil {
				return 0, idx, fmt.Errorf("failed to unpack subfield %s (%s): %w", idx, fl.Spec().Description, err)
			}

			offset += read
		}
	}

	return offset, "", nil
}

const (
	// ignoredMaxLen is a constant meant to be used in encoders that don't use the maxLength argument during
	// length decoding.
	ignoredMaxLen int = 0
	// maxLenOfUnknownTag is max int in order to never hit this limit.
	maxLenOfUnknownTag = math.MaxInt
)

func (f *Composite) unpackSubfieldsByTag(data []byte) (int, string, error) {
	// reset subfields
	f.subfields = make(map[string]Field)

	offset := 0

	for offset < len(data) {
		tagBytes, read, err := f.spec.Tag.Enc.Decode(data[offset:], f.spec.Tag.Length)
		if err != nil {
			return 0, "", fmt.Errorf("failed to unpack subfield Tag: %w", err)
		}

		offset += read

		if f.spec.Tag.Pad != nil {
			tagBytes = f.spec.Tag.Pad.Unpad(tagBytes)
		}

		tag := string(tagBytes)

		specField, ok := f.spec.Subfields[tag]
		if !ok {
			if f.skipUnknownTLVTags() {
				// to obtain the length of the unknown tag and add it to the offset we need to decode its length
				// by default, we use BER-TVL prefix because BER-TLV lengths are decoded dynamically, the maxLen method argument is ignored.
				var (
					pref   prefix.Prefixer = prefix.BerTLV
					maxLen                 = ignoredMaxLen
				)

				// but if PrefUnknownTLV prefix is set, use it and hope that length of all unknown tags is encoded using this prefixer
				if f.spec.Tag.PrefUnknownTLV != nil {
					pref = f.spec.Tag.PrefUnknownTLV
					maxLen = maxLenOfUnknownTag
				}

				fieldLength, read, err := pref.DecodeLength(maxLen, data[offset:])
				if err != nil {
					return 0, "", err
				}

				// store unknown field as Binary if StoreUnknownTLVTags is enabled
				if f.spec.Tag.StoreUnknownTLVTags {
					fieldData := data[offset+read : offset+read+fieldLength]
					binaryField := NewBinary(&Spec{
						Length:      fieldLength,
						Description: fmt.Sprintf("Unknown TLV tag %s", tag),
						Pref:        pref,
						Enc:         encoding.Binary,
					})
					if err := binaryField.SetBytes(fieldData); err != nil {
						return 0, tag, fmt.Errorf("failed to set bytes for unknown tag %s: %w", tag, err)
					}
					f.subfields[tag] = binaryField
				}

				offset += fieldLength + read

				continue
			}

			return 0, tag, fmt.Errorf("failed to unpack subfield %v: field is not defined in the spec", tag)
		}

		field := NewInstanceOf(specField)
		f.subfields[tag] = field

		read, err = field.Unpack(data[offset:])
		if err != nil {
			return 0, tag, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}

		offset += read
	}

	return offset, "", nil
}

func (f *Composite) skipUnknownTLVTags() bool {
	return f.spec.Tag != nil && f.spec.Tag.SkipUnknownTLVTags && (f.spec.Tag.Enc == encoding.BerTLVTag || f.spec.Tag.PrefUnknownTLV != nil)
}

func orderedKeys(kvs map[string]Field, sorter sort.StringSlice) []string {
	keys := make([]string, 0, len(kvs))
	for k := range kvs {
		keys = append(keys, k)
	}
	sorter(keys)

	return keys
}

// UnsetSubfield marks the subfield with the given ID as not set and replaces it
// with a new zero-valued field. This effectively removes the subfield's value and
// excludes it from operations like Pack() or Marshal().
func (m *Composite) UnsetSubfield(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.unsetField(id)
}

// isn't protected by mutex, caller must ensure mutex is locked
func (m *Composite) unsetField(id string) {
	delete(m.subfields, id)
}

// UnsetSubfields marks multiple subfields identified by their paths as not set and
// replaces them with new zero-valued fields. Each path should be in the format
// "a.b.c". This effectively removes the subfields' values and excludes them from
// operations like Pack() or Marshal().
// Deprecated: use UnsetPath instead.
func (m *Composite) UnsetSubfields(idPaths ...string) error {
	return m.UnsetPath(idPaths...)
}

// UnsetPath marks field identified by their path as not set and replaces them
// with new zero-valued fields. Each path should be in the format "a.b.c". It
// accepts multiple paths as arguments.
func (m *Composite) UnsetPath(idPaths ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, idPath := range idPaths {
		if idPath == "" {
			continue
		}

		id, path, hasSubPath := strings.Cut(idPath, ".")

		f := m.subfields[id]
		if f == nil {
			continue
		}

		if !hasSubPath {
			m.unsetField(id)
			continue
		}

		pathUnsetter, ok := f.(PathUnsetter)
		if !ok {
			return fmt.Errorf("field %s is not a composite field and its subfields %s cannot be unset", id, path)
		}

		if err := pathUnsetter.UnsetPath(path); err != nil {
			return fmt.Errorf("failed to unset %s in composite field %s: %w", path, id, err)
		}
	}

	return nil
}

func (m *Composite) MarshalPath(path string, value any) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id, subPath, hasSubPath := strings.Cut(path, ".")

	field, err := m.getOrCreateField(id)
	if err != nil {
		return fmt.Errorf("getting or creating subfield %s: %w", id, err)
	}

	// if there is subPath, marshal it recursively
	if hasSubPath {
		// check if field supports MarshalPath
		mField, ok := field.(PathMarshaler)
		if !ok {
			return fmt.Errorf("field %s is not a PathMarshaler", id)
		}

		err := mField.MarshalPath(subPath, value)
		if err != nil {
			return fmt.Errorf("marshaling path %s in field %s: %w", subPath, id, err)
		}

		return nil
	}

	err = field.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling field %s: %w", id, err)
	}

	return nil
}

func (m *Composite) UnmarshalPath(path string, value any) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id, subPath, hasSubPath := strings.Cut(path, ".")

	field := m.subfields[id]
	if field == nil {
		if _, ok := m.Spec().Subfields[id]; !ok {
			return fmt.Errorf("field %s is not defined in the spec", id)
		}

		return nil
	}

	// if there is subPath, unmarshal it recursively
	if hasSubPath {
		// check if field supports UnmarshalPath
		uField, ok := field.(PathUnmarshaler)
		if !ok {
			return fmt.Errorf("field %s is not a PathUnmarshaler", id)
		}

		err := uField.UnmarshalPath(subPath, value)
		if err != nil {
			return fmt.Errorf("unmarshaling path %s in field %s: %w", subPath, id, err)
		}

		return nil
	}

	err := field.Unmarshal(value)
	if err != nil {
		return fmt.Errorf("unmarshaling field %s: %w", id, err)
	}

	return nil
}
