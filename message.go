package iso8583

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/moov-io/iso8583/field"
	"github.com/stoewer/go-strcase"
)

const (
	fieldStringTypeLabel  = "*field.String"
	fieldBitmapTypeLabel  = "*field.Bitmap"
	fieldNumericTypeLabel = "*field.Numeric"
	messageLabel          = "ISO8583"
)

type dummyMap map[string]string

func (m dummyMap) sort() (index []string) {
	for k := range m {
		index = append(index, k)
	}
	sort.Strings(index)
	return
}

type Message struct {
	fields    map[int]field.Field
	spec      *MessageSpec
	data      interface{}
	fieldsMap map[int]struct{}
	bitmap    *field.Bitmap
}

func NewMessage(spec *MessageSpec) *Message {
	fields := spec.CreateMessageFields()

	return &Message{
		fields:    fields,
		spec:      spec,
		fieldsMap: map[int]struct{}{},
	}
}

func (m *Message) Data() interface{} {
	return m.data
}

func (m *Message) SetData(data interface{}) error {
	m.data = data

	if m.data == nil {
		return nil
	}

	// get the struct
	str := reflect.ValueOf(data).Elem()

	if reflect.TypeOf(str).Kind() != reflect.Struct {
		return fmt.Errorf("failed to set data as struct is expected, got: %s", reflect.TypeOf(str).Kind())
	}

	for i, fl := range m.fields {
		fieldName := fmt.Sprintf("F%d", i)

		// get the struct field
		dataField := str.FieldByName(fieldName)

		if dataField == (reflect.Value{}) || dataField.IsNil() {
			continue
		}

		if dataField.Type() != reflect.TypeOf(fl) {
			return fmt.Errorf("failed to set data: type of the field %d: %v does not match the type in the spec: %v", i, dataField.Type(), reflect.TypeOf(fl))
		}

		// set data field spec for the message spec field
		specField := m.fields[i]
		df := dataField.Interface().(field.Field)
		df.SetSpec(specField.Spec())

		// use data field as a message field
		m.fields[i] = df
		m.fieldsMap[i] = struct{}{}
	}

	return nil
}

func (m *Message) Bitmap() *field.Bitmap {
	if m.bitmap != nil {
		return m.bitmap
	}

	m.bitmap = m.fields[1].(*field.Bitmap)
	m.fieldsMap[1] = struct{}{}

	return m.bitmap
}

func (m *Message) MTI(val string) {
	m.fieldsMap[0] = struct{}{}
	m.fields[0].SetBytes([]byte(val))
}

func (m *Message) Field(id int, val string) {
	m.fieldsMap[id] = struct{}{}
	m.fields[id].SetBytes([]byte(val))
}

func (m *Message) BinaryField(id int, val []byte) {
	m.fieldsMap[id] = struct{}{}
	m.fields[id].SetBytes(val)
}

func (m *Message) GetMTI() string {
	// check index
	return m.fields[0].String()
}

func (m *Message) GetString(id int) string {
	if _, ok := m.fieldsMap[id]; ok {
		return m.fields[id].String()
	}

	return ""
}

func (m *Message) GetBytes(id int) []byte {
	if _, ok := m.fieldsMap[id]; ok {
		return m.fields[id].Bytes()
	}

	return nil
}

func (m *Message) Pack() ([]byte, error) {

	packed := make([]byte, 0)
	m.Bitmap().Reset()

	// build the bitmap
	var fieldIndexes []int
	for id := range m.fieldsMap {

		fieldIndexes = append(fieldIndexes, id)

		// indexes 0 and 1 are for mti and bitmap
		// regular field number started from index 2
		if id < 2 {
			continue
		}

		m.Bitmap().Set(id)
	}

	// pack fields
	sort.Ints(fieldIndexes)
	for _, index := range fieldIndexes {
		field, ok := m.fields[index]
		if !ok {
			return nil, fmt.Errorf("failed to pack field %d: no specification found", index)
		}

		packedField, err := field.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to pack field %d (%s): %v", index, field.Spec().Description, err)
		}
		packed = append(packed, packedField...)
	}

	return packed, nil
}

func (m *Message) Unpack(src []byte) error {
	var off int

	m.fieldsMap = map[int]struct{}{}
	m.Bitmap().Reset()

	// unpack MTI
	read, err := m.fields[0].Unpack(src)
	if err != nil {
		return err
	}

	off = read

	// unpack Bitmap
	read, err = m.fields[1].Unpack(src[off:])
	if err != nil {
		return err
	}

	off += read

	for i := 2; i <= m.Bitmap().Len(); i++ {
		if m.Bitmap().IsSet(i) {
			field, ok := m.fields[i]
			if !ok {
				return fmt.Errorf("failed to unpack field %d: no specification found", i)
			}

			m.fieldsMap[i] = struct{}{}
			read, err = field.Unpack(src[off:])
			if err != nil {
				return fmt.Errorf("failed to unpack field %d (%s): %v", i, field.Spec().Description, err)
			}

			err = m.linkDataFieldWithMessageField(i, field)
			if err != nil {
				return fmt.Errorf("failed to unpack field %d: %v", i, err)
			}
			off += read
		}
	}

	return nil
}

// Customize unmarshal of json
func (m *Message) UnmarshalJSON(b []byte) error {

	if m.spec == nil {
		return errors.New("please set specification of message")
	}

	dummy := dummyMap{}

	err := json.Unmarshal(b, &dummy)
	if err != nil {
		return errors.New("failed to parse json string of message")
	}

	return m.setDataWithMap(dummy, "SnakeCase")
}

// Customize marshal of json
func (m *Message) MarshalJSON() ([]byte, error) {

	if m.spec == nil {
		return nil, errors.New("please set specification of message")
	}

	// after encoding
	m.Pack()

	var fieldIndexes []int
	for id := range m.fieldsMap {
		fieldIndexes = append(fieldIndexes, id)
	}

	dummy := dummyMap{}

	sort.Ints(fieldIndexes)
	for _, index := range fieldIndexes {
		field, ok := m.fields[index]
		if !ok {
			return nil, fmt.Errorf("failed to pack field %d: no specification found", index)
		}

		fieldName := strcase.SnakeCase(field.Spec().GetIdentifier())
		if len(fieldName) > 0 {
			dummy[fieldName] = field.String()
		}

	}

	dummy.sort()
	return json.Marshal(dummy)
}

// Customize marshal of xml
func (m *Message) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

	if m.spec == nil {
		return errors.New("please set specification of message")
	}

	dummy := dummyMap{}
	values := make([]string, 0)

	for token, err := decoder.Token(); err == nil; token, err = decoder.Token() {
		switch t := token.(type) {
		case xml.CharData:
			values = append(values, string([]byte(t)))
		case xml.EndElement:
			if t.Name.Local == "langs" || t.Name.Local == messageLabel {
				continue
			}

			var value string
			err = json.Unmarshal([]byte("\""+values[len(values)-1]+"\""), &value)
			if err != nil {
				continue
			}

			dummy[t.Name.Local] = value
		}
	}

	return m.setDataWithMap(dummy, "CamelCase")
}

// Customize unmarshal of xml
func (m *Message) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {

	if m.spec == nil {
		return errors.New("please set specification of message")
	}

	// after encoding
	m.Pack()

	var fieldIndexes []int
	for id := range m.fieldsMap {
		fieldIndexes = append(fieldIndexes, id)
	}

	start.Name = xml.Name{Local: messageLabel}
	tokens := []xml.Token{start}

	sort.Ints(fieldIndexes)
	for _, index := range fieldIndexes {
		field, ok := m.fields[index]
		if !ok {
			return fmt.Errorf("failed to marshal field %d: no specification found", index)
		}

		t := xml.StartElement{
			Name: xml.Name{Local: strcase.UpperCamelCase(field.Spec().GetIdentifier())},
		}

		val, err := json.Marshal(field.String())
		if err != nil {
			return err
		}

		tokens = append(tokens, t, xml.CharData(strings.Trim(string(val), "\"")), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{Name: start.Name})
	for _, t := range tokens {
		err := encoder.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	return encoder.Flush()
}

func (m *Message) linkDataFieldWithMessageField(i int, fl field.Field) error {
	if m.data == nil {
		return nil
	}

	// get the struct
	str := reflect.ValueOf(m.data).Elem()

	fieldName := fmt.Sprintf("F%d", i)

	// get the struct field
	dataField := str.FieldByName(fieldName)
	if dataField == (reflect.Value{}) {
		return nil
	}

	if dataField.Type() != reflect.TypeOf(fl) {
		return fmt.Errorf("field type: %v does not match the type in the spec: %v", dataField.Type(), reflect.TypeOf(fl))
	}

	dataField.Addr().Elem().Set(reflect.ValueOf(fl))

	return nil
}

func (m *Message) getType(field interface{}) string {
	return reflect.TypeOf(field).String()
}

func (m *Message) createNewFieldWithValue(f field.Field, value string) (newField field.Field, err error) {
	wantFieldName := m.getType(f)
	switch wantFieldName {
	case fieldStringTypeLabel:
		newField = field.NewStringValue(value)
		newField.SetSpec(f.Spec())
	case fieldBitmapTypeLabel:
		newField = field.NewBitmap(f.Spec())
		newField.SetBytes([]byte(value))
	case fieldNumericTypeLabel:
		var intValue int
		intValue, err = strconv.Atoi(value)
		if err != nil {
			return
		}
		newField = field.NewNumericValue(intValue)
		newField.SetSpec(f.Spec())
	default:
		err = errors.New("has unsupported field type")
		return
	}

	return
}

func (m *Message) setDataWithMap(dummy map[string]string, stringCase string) error {
	for key, value := range dummy {

		index, err := m.spec.GetFieldIndex(key, stringCase)
		if err != nil {
			fmt.Println(err)
			continue
		}

		specField := m.spec.Fields[index]
		f, err := m.createNewFieldWithValue(specField, value)
		if err != nil {
			return err
		}

		m.fields[index] = f
		m.fieldsMap[index] = struct{}{}
	}
	return nil
}
