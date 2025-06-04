package field

import (
	"encoding/json"
	"errors"
	"fmt"
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

// ConstructSubfields creates subfields according to the spec
// this method is used when composite field is created without
// calling NewComposite (when we create message spec and composite spec)
func (f *Composite) ConstructSubfields() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.subfields == nil {
		f.subfields = CreateSubfields(f.spec)
	}
	f.setSubfields = make(map[string]struct{})
}

// Spec returns the receiver's spec.
func (f *Composite) Spec() *Spec {
	return f.spec
}

// GetSubfields returns the map of set sub fields
func (f *Composite) GetSubfields() map[string]Field {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.getSubfields()
}

// getSubfields returns the map of set sub fields, it should be called
// only when the mutex is locked
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
// should only pass None or nil values for this type. Passing any other value
// will result in a panic.
func (f *Composite) SetSpec(spec *Spec) {
	if err := spec.Validate(); err != nil {
		panic(err) //nolint:forbidigo,nolintlint // as specs mostly static, we panic on spec validation errors
	}
	f.spec = spec

	var sortFn sort.StringSlice

	// When bitmap is not defined, always order tags by int.
	if spec.Bitmap != nil {
		sortFn = sort.StringsByInt
	} else {
		sortFn = spec.Tag.Sort
	}

	f.orderedSpecFieldTags = orderedKeys(spec.Subfields, sortFn)
}

func (f *Composite) Unmarshal(v interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

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
		indexTag := NewIndexTag(dataStruct.Type().Field(i))
		if indexTag.Tag == "" {
			continue
		}

		messageField, ok := f.subfields[indexTag.Tag]
		if !ok {
			continue
		}

		// unmarshal only subfield that has the value set
		if _, set := f.setSubfields[indexTag.Tag]; !set {
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
	f.mu.Lock()
	defer f.mu.Unlock()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
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
	for i := 0; i < dataStruct.NumField(); i++ {
		indexTag := NewIndexTag(dataStruct.Type().Field(i))
		if indexTag.Tag == "" {
			continue
		}

		messageField, ok := f.subfields[indexTag.Tag]
		if !ok {
			continue
		}

		dataField := dataStruct.Field(i)
		if dataField.IsZero() && !indexTag.KeepZero {
			continue
		}

		err := messageField.Marshal(dataField.Interface())
		if err != nil {
			return fmt.Errorf("marshalling field %s: %w", indexTag.Tag, err)
		}

		f.setSubfields[indexTag.Tag] = struct{}{}
	}

	return nil
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

	bitmap, ok := CreateSubfield(f.spec.Bitmap).(*Bitmap)
	if !ok {
		return nil
	}

	f.cachedBitmap = bitmap

	return f.cachedBitmap
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
	if f.bitmap() != nil {
		return f.packByBitmap()
	}

	return f.packByTag()
}

func (f *Composite) packByBitmap() ([]byte, error) {
	f.bitmap().Reset()

	var packedFields []byte

	// pack fields
	for _, id := range f.orderedSpecFieldTags {
		// If this ordered field is not set, continue to the next field.
		if _, ok := f.setSubfields[id]; !ok {
			continue
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("converting id %s to int: %w", id, err)
		}

		// set bitmap bit for this field
		f.bitmap().Set(idInt)

		field, ok := f.subfields[id]
		if !ok {
			return nil, fmt.Errorf("failed to pack subfield %s: no specification found", id)
		}

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
	if f.bitmap() != nil {
		return f.unpackSubfieldsByBitmap(data)
	}
	if f.spec.Tag.Enc != nil {
		return f.unpackSubfieldsByTag(data)
	}
	return f.unpackSubfields(data, isVariableLength)
}

func (f *Composite) unpackSubfields(data []byte, isVariableLength bool) (int, string, error) {
	offset := 0
	for _, tag := range f.orderedSpecFieldTags {
		field, ok := f.subfields[tag]
		if !ok {
			continue
		}

		read, err := field.Unpack(data[offset:])
		if err != nil {
			return 0, tag, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}

		f.setSubfields[tag] = struct{}{}

		offset += read

		if isVariableLength && offset >= len(data) {
			break
		}
	}

	return offset, "", nil
}

func (f *Composite) unpackSubfieldsByBitmap(data []byte) (int, string, error) {
	var off int

	// Reset fields that were set.
	f.setSubfields = make(map[string]struct{})

	f.bitmap().Reset()

	read, err := f.bitmap().Unpack(data[off:])
	if err != nil {
		return 0, "", fmt.Errorf("failed to unpack bitmap: %w", err)
	}

	off += read

	for i := 1; i <= f.bitmap().Len(); i++ {
		if f.bitmap().IsSet(i) {
			iStr := strconv.Itoa(i)

			fl, ok := f.subfields[iStr]
			if !ok {
				return 0, iStr, fmt.Errorf("failed to unpack subfield %s: no specification found", iStr)
			}

			read, err = fl.Unpack(data[off:])
			if err != nil {
				return 0, iStr, fmt.Errorf("failed to unpack subfield %s (%s): %w", iStr, fl.Spec().Description, err)
			}

			f.setSubfields[iStr] = struct{}{}

			off += read
		}
	}

	return off, "", nil
}

const (
	// ignoredMaxLen is a constant meant to be used in encoders that don't use the maxLength argument during
	// length decoding.
	ignoredMaxLen int = 0
	// maxLenOfUnknownTag is max int in order to never hit this limit.
	maxLenOfUnknownTag = math.MaxInt
)

func (f *Composite) unpackSubfieldsByTag(data []byte) (int, string, error) {
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
		if _, ok := f.spec.Subfields[tag]; !ok {
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
				offset += fieldLength + read
				continue
			}

			return 0, tag, fmt.Errorf("failed to unpack subfield %v: field not defined in Spec", tag)
		}

		field, ok := f.subfields[tag]
		if !ok {
			continue
		}

		read, err = field.Unpack(data[offset:])
		if err != nil {
			return 0, tag, fmt.Errorf("failed to unpack subfield %v: %w", tag, err)
		}

		f.setSubfields[tag] = struct{}{}

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

	// unset the field
	delete(m.setSubfields, id)

	// we should re-create the subfield to reset its value (and its subfields)
	m.subfields[id] = CreateSubfield(m.Spec().Subfields[id])
}

// UnsetSubfields marks multiple subfields identified by their paths as not set and
// replaces them with new zero-valued fields. Each path should be in the format
// "a.b.c". This effectively removes the subfields' values and excludes them from
// operations like Pack() or Marshal().
func (m *Composite) UnsetSubfields(idPaths ...string) error {
	for _, idPath := range idPaths {
		if idPath == "" {
			continue
		}

		id, path, _ := strings.Cut(idPath, ".")

		if _, ok := m.setSubfields[id]; ok {
			if len(path) == 0 {
				m.UnsetSubfield(id)
				continue
			}

			f := m.subfields[id]
			if f == nil {
				return fmt.Errorf("subfield %s does not exist", id)
			}

			composite, ok := f.(*Composite)
			if !ok {
				return fmt.Errorf("field %s is not a composite field and its subfields %s cannot be unset", id, path)
			}

			if err := composite.UnsetSubfields(path); err != nil {
				return fmt.Errorf("failed to unset %s in composite field %s: %w", path, id, err)
			}
		}
	}

	return nil
}
