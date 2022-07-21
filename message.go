package iso8583

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"

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
	spec   *MessageSpec
	bitmap *field.Bitmap

	// stores all fields according to the spec
	fields map[int]field.Field

	// tracks which fields were set
	fieldsMap map[int]struct{}
}

func NewMessage(spec *MessageSpec) *Message {
	fields := spec.CreateMessageFields()

	return &Message{
		fields:    fields,
		spec:      spec,
		fieldsMap: map[int]struct{}{},
	}
}

// Deprecated. Use Marshal intead.
func (m *Message) SetData(data interface{}) error {
	return m.Marshal(data)
}

func (m *Message) Bitmap() *field.Bitmap {
	if m.bitmap != nil {
		return m.bitmap
	}

	m.bitmap = m.fields[bitmapIdx].(*field.Bitmap)
	m.bitmap.Reset()
	m.fieldsMap[bitmapIdx] = struct{}{}

	return m.bitmap
}

func (m *Message) MTI(val string) {
	m.fieldsMap[mtiIdx] = struct{}{}
	m.fields[mtiIdx].SetBytes([]byte(val))
}

func (m *Message) GetSpec() *MessageSpec {
	return m.spec
}

func (m *Message) Field(id int, val string) error {
	if f, ok := m.fields[id]; ok {
		m.fieldsMap[id] = struct{}{}
		return f.SetBytes([]byte(val))
	}
	return fmt.Errorf("failed to set field %d. ID does not exist", id)
}

func (m *Message) BinaryField(id int, val []byte) error {
	if f, ok := m.fields[id]; ok {
		m.fieldsMap[id] = struct{}{}
		return f.SetBytes(val)
	}
	return fmt.Errorf("failed to set binary field %d. ID does not exist", id)
}

func (m *Message) GetMTI() (string, error) {
	// check index
	return m.fields[mtiIdx].String()
}

func (m *Message) GetString(id int) (string, error) {
	if f, ok := m.fields[id]; ok {
		m.fieldsMap[id] = struct{}{}
		return f.String()
	}
	return "", fmt.Errorf("failed to get string for field %d. ID does not exist", id)
}

func (m *Message) GetBytes(id int) ([]byte, error) {
	if f, ok := m.fields[id]; ok {
		m.fieldsMap[id] = struct{}{}
		return f.Bytes()
	}
	return nil, fmt.Errorf("failed to get bytes for field %d. ID does not exist", id)
}

func (m *Message) GetField(id int) field.Field {
	return m.fields[id]
}

// Fields returns the map of the set fields
func (m *Message) GetFields() map[int]field.Field {
	fields := map[int]field.Field{}
	for i := range m.fieldsMap {
		fields[i] = m.GetField(i)
	}
	return fields
}

func (m *Message) Pack() ([]byte, error) {
	packed := []byte{}
	m.Bitmap().Reset()

	ids, err := m.packableFieldIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to pack message: %w", err)
	}

	for _, id := range ids {
		// indexes 0 and 1 are for mti and bitmap
		// regular field number startd from index 2
		if id < 2 {
			continue
		}
		m.Bitmap().Set(id)
	}

	// pack fields
	for _, i := range ids {
		field, ok := m.fields[i]
		if !ok {
			return nil, fmt.Errorf("failed to pack field %d: no specification found", i)
		}
		packedField, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack field %d (%s): %w", i, field.Spec().Description, err)
		}
		packed = append(packed, packedField...)
	}

	return packed, nil
}

func (m *Message) Unpack(src []byte) error {
	var off int

	// reset fields that were set
	m.fieldsMap = map[int]struct{}{}

	// This method implicitly also sets m.fieldsMap[bitmapIdx]
	m.Bitmap().Reset()

	read, err := m.fields[mtiIdx].Unpack(src)
	if err != nil {
		return fmt.Errorf("failed to unpack MTI: %w", err)
	}

	m.fieldsMap[mtiIdx] = struct{}{}

	off = read

	// unpack Bitmap
	read, err = m.fields[bitmapIdx].Unpack(src[off:])
	if err != nil {
		return fmt.Errorf("failed to unpack bitmap: %w", err)
	}

	off += read

	for i := 2; i <= m.Bitmap().Len(); i++ {
		if m.Bitmap().IsSet(i) {
			fl, ok := m.fields[i]
			if !ok {
				return fmt.Errorf("failed to unpack field %d: no specification found", i)
			}

			read, err = fl.Unpack(src[off:])
			if err != nil {
				return fmt.Errorf("failed to unpack field %d (%s): %w", i, fl.Spec().Description, err)
			}

			m.fieldsMap[i] = struct{}{}

			off += read
		}
	}

	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	// by packing the message we will generate bitmap
	// create HEX representation
	// and validate message against the spec
	if _, err := m.Pack(); err != nil {
		return nil, err
	}

	fieldMap := m.GetFields()
	strFieldMap := map[string]field.Field{}
	for id, field := range fieldMap {
		// we don't wish to populate the bitmap in the final
		// JSON since it is dynamically generated when packing
		// and unpacking anyways.
		if id == bitmapIdx {
			continue
		}
		strFieldMap[fmt.Sprint(id)] = field
	}

	// get only fields that were set
	bytes, err := json.Marshal(field.OrderedMap(strFieldMap))
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal map to bytes")
	}
	return bytes, nil
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for id, rawMsg := range data {
		i, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field %v: could not convert to int", i)
		}

		field, ok := m.fields[i]
		if !ok {
			return fmt.Errorf("failed to unmarshal field %d: no specification found", i)
		}

		if err := json.Unmarshal(rawMsg, field); err != nil {
			return utils.NewSafeErrorf(err, "failed to unmarshal field %v", id)
		}

		m.fieldsMap[i] = struct{}{}
	}

	return nil
}

func (m *Message) packableFieldIDs() ([]int, error) {
	// Index 1 represent bitmap which is always populated.
	populatedFieldIDs := []int{1}

	for id := range m.fieldsMap {
		// represents the bitmap
		if id == 1 {
			continue
		}

		populatedFieldIDs = append(populatedFieldIDs, id)
	}

	sort.Ints(populatedFieldIDs)

	return populatedFieldIDs, nil
}

// Clone clones the message by creating a new message from the binary
// representation of the original message
func (m *Message) Clone() (*Message, error) {
	newMessage := NewMessage(m.spec)

	bytes, err := m.Pack()
	if err != nil {
		return nil, err
	}

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

	// iterate over struct fields
	for i := 0; i < dataStruct.NumField(); i++ {
		fieldIndex, err := getFieldIndex(dataStruct.Type().Field(i))
		if err != nil {
			return fmt.Errorf("getting field %d index: %w", i, err)
		}

		// skip field without index
		if fieldIndex < 0 {
			continue
		}

		messageField := m.GetField(fieldIndex)
		// if struct field we are usgin to populate value expects to
		// set index of the field that is not described by spec
		if messageField == nil {
			return fmt.Errorf("no message field defined by spec with index: %d", fieldIndex)
		}

		dataField := dataStruct.Field(i)
		if dataField.IsNil() {
			continue
		}

		err = messageField.Marshal(dataField.Interface())
		if err != nil {
			return fmt.Errorf("failed to set value to field %d: %w", fieldIndex, err)
		}

		m.fieldsMap[fieldIndex] = struct{}{}
	}

	return nil
}

// Unmarshal populates v struct fields with message field values. It traverses
// through the message fields and calls Unmarshal(...) on them setting the v If
// v  is nil or not a pointer it returns error.
func (m *Message) Unmarshal(v interface{}) error {
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
		fieldIndex, err := getFieldIndex(dataStruct.Type().Field(i))
		if err != nil {
			return fmt.Errorf("getting field %d index: %w", i, err)
		}

		// skip field without index
		if fieldIndex < 0 {
			continue
		}

		// we can get data only if field value is set
		messageField := m.GetField(fieldIndex)
		if messageField == nil {
			continue
		}

		if _, set := m.fieldsMap[fieldIndex]; !set {
			continue
		}

		dataField := dataStruct.Field(i)
		if dataField.IsNil() {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}

		err = messageField.Unmarshal(dataField.Interface())
		if err != nil {
			return fmt.Errorf("failed to get value from field %d: %w", fieldIndex, err)
		}
	}

	return nil
}

var fieldNameIndexRe = regexp.MustCompile(`^F\d+$`)

// fieldIndex returns index of the field. First, it checks field name. If it
// does not match FNN (when NN is digits), it checks value of `index` tag.  If
// negative value returned (-1) then index was not found for the field.
func getFieldIndex(field reflect.StructField) (int, error) {
	dataFieldName := field.Name

	if indexStr := field.Tag.Get("index"); indexStr != "" {
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return -1, fmt.Errorf("converting field index into int: %w", err)
		}

		return fieldIndex, nil
	}

	if len(dataFieldName) > 0 && fieldNameIndexRe.MatchString(dataFieldName) {
		indexStr := dataFieldName[1:]
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return -1, fmt.Errorf("converting field index into int: %w", err)
		}

		return fieldIndex, nil
	}

	return -1, nil
}
