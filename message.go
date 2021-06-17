package iso8583

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/field"
)

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

	dataStruct := reflect.ValueOf(data)
	strKind := dataStruct.Kind()
	if strKind == reflect.Ptr || strKind == reflect.Interface {
		// get the struct
		dataStruct = dataStruct.Elem()
	}

	if reflect.TypeOf(dataStruct).Kind() != reflect.Struct {
		return fmt.Errorf("failed to set data as struct is expected, got: %s", reflect.TypeOf(dataStruct).Kind())
	}

	for i, specField := range m.fields {
		fieldName := fmt.Sprintf("F%d", i)

		// get the struct field
		dataField := dataStruct.FieldByName(fieldName)

		// no data field was found with this name
		if dataField == (reflect.Value{}) {
			continue
		}

		isNil := dataField.IsNil()
		if isNil {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}
		if err := specField.SetData(dataField.Interface()); err != nil {
			return fmt.Errorf("failed to set data for field %d: %w", i, err)
		}
		if !isNil {
			m.fieldsMap[i] = struct{}{}
		}

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
	return m.fields[0].String()
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

func (m *Message) Pack() ([]byte, error) {
	packed := []byte{}
	m.Bitmap().Reset()

	// build the bitmap
	maxId := 0

	for id := range m.fieldsMap {
		if id > maxId {
			maxId = id
		}

		// indexes 0 and 1 are for mti and bitmap
		// regular field number startd from index 2
		if id < 2 {
			continue
		}

		m.Bitmap().Set(id)
	}

	// pack fields
	for i := 0; i <= maxId; i++ {
		if _, ok := m.fieldsMap[i]; ok {
			field, ok := m.fields[i]
			if !ok {
				return nil, fmt.Errorf("failed to pack field %d: no specification found", i)
			}

			packedField, err := field.Pack()
			if err != nil {
				return nil, fmt.Errorf("failed to pack field %d (%s): %v", i, field.Spec().Description, err)
			}
			packed = append(packed, packedField...)
		}
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
		return fmt.Errorf("failed to unpack MTI: %v", err)
	}

	off = read

	// unpack Bitmap
	read, err = m.fields[1].Unpack(src[off:])
	if err != nil {
		return fmt.Errorf("failed to unpack bitmapt: %v", err)
	}

	off += read

	for i := 2; i <= m.Bitmap().Len(); i++ {
		if m.Bitmap().IsSet(i) {
			fl, ok := m.fields[i]
			if !ok {
				return fmt.Errorf("failed to unpack field %d: no specification found", i)
			}

			m.fieldsMap[i] = struct{}{}
			read, err = fl.Unpack(src[off:])
			if err != nil {
				return fmt.Errorf("failed to unpack field %d (%s): %v", i, fl.Spec().Description, err)
			}

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

	jsonData := field.OrderedMap(m.fields)

	return json.Marshal(jsonData)
}
