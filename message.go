package iso8583

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/moov-io/iso8583/field"
)

type Message struct {
	fields    map[int]field.Field
	spec      *MessageSpec
	data      interface{}
	dataValue *reflect.Value
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
	if dataStruct.Kind() == reflect.Ptr || dataStruct.Kind() == reflect.Interface {
		// get the struct
		dataStruct = dataStruct.Elem()
	}

	if reflect.TypeOf(dataStruct).Kind() != reflect.Struct {
		return fmt.Errorf("failed to set data as struct is expected, got: %s", reflect.TypeOf(dataStruct).Kind())
	}

	m.dataValue = &dataStruct
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

func (m *Message) GetField(id int) (field.Field, error) {
	if f, ok := m.fields[id]; ok {
		return f, nil
	}

	return nil, fmt.Errorf("failed to get the field %d. ID does not exist", id)
}

func (m *Message) Pack() ([]byte, error) {
	var buf bytes.Buffer

	_, err := m.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (m *Message) Unpack(src []byte) error {
	_, err := m.ReadFrom(bytes.NewReader(src))
	return err
}

func (m *Message) WriteTo(w io.Writer) (n int, err error) {
	m.Bitmap().Reset()

	ids, err := m.setPackableDataFields()
	if err != nil {
		return 0, fmt.Errorf("failed to pack message: %w", err)
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
		field, err := m.GetField(i)
		if err != nil {
			return 0, err
		}

		m, err := field.WriteTo(w)
		if err != nil {
			return 0, fmt.Errorf("failed to pack field %d (%s): %v", i, field.Spec().Description, err)
		}
		n += m
	}

	return n, nil
}

func (m *Message) ReadFrom(r io.Reader) (n int, err error) {
	m.fieldsMap = map[int]struct{}{}
	m.Bitmap().Reset()

	// unpack MTI
	mti, err := m.GetField(0)
	if err != nil {
		return 0, err
	}

	read, err := mti.ReadFrom(r)
	if err != nil {
		return read, fmt.Errorf("failed to unpack MTI: %v", err)
	}
	n += read

	// unpack Bitmap
	bitmap, err := m.GetField(1)
	if err != nil {
		return n, err
	}

	read, err = bitmap.ReadFrom(r)
	if err != nil {
		return n + read, fmt.Errorf("failed to unpack bitmapt: %v", err)
	}
	n += read

	for i := 2; i <= m.Bitmap().Len(); i++ {
		if m.Bitmap().IsSet(i) {
			fl, err := m.GetField(i)
			if err != nil {
				return n, err
			}

			if m.dataValue != nil {
				if err := m.setUnpackableDataField(i, fl); err != nil {
					return n, err
				}
			}

			m.fieldsMap[i] = struct{}{}
			read, err = fl.ReadFrom(r)
			if err != nil {
				return n + read, fmt.Errorf("failed to unpack field %d (%s): %v", i, fl.Spec().Description, err)
			}
			n += read
		}
	}

	return n, nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	// by packing the message we will generate bitmap
	// create HEX representation
	// and validate message against the spec
	if _, err := m.WriteTo(&bytes.Buffer{}); err != nil {
		return nil, err
	}

	jsonData := field.OrderedMap(m.fields)

	return json.Marshal(jsonData)
}

func (m *Message) setPackableDataFields() ([]int, error) {
	// Indexes 0 and 1 represent the mti and bitmap.
	// These fields are always populated.
	populatedFieldIDs := []int{0, 1}

	for id, field := range m.fields {
		// regular field number start from index 2
		if id < 2 {
			continue
		}

		// These fields are set using the typed API
		if m.dataValue != nil {
			dataField := m.dataFieldValue(id)
			// no non-nil data field was found with this name
			if dataField == (reflect.Value{}) || dataField.IsNil() {
				continue
			}
			if err := field.SetData(dataField.Interface()); err != nil {
				return nil, fmt.Errorf("failed to set data for field %d: %w", id, err)
			}
		}

		// These fields are set using the untyped API
		_, ok := m.fieldsMap[id]
		if ok || m.dataValue != nil {
			populatedFieldIDs = append(populatedFieldIDs, id)
		}
	}
	sort.Ints(populatedFieldIDs)

	return populatedFieldIDs, nil
}

func (m *Message) setUnpackableDataField(id int, specField field.Field) error {
	dataField := m.dataFieldValue(id)
	// no data field was found with this name
	if dataField == (reflect.Value{}) {
		return nil
	}

	isNil := dataField.IsNil()
	if isNil {
		dataField.Set(reflect.New(dataField.Type().Elem()))
	}
	if err := specField.SetData(dataField.Interface()); err != nil {
		return fmt.Errorf("failed to set data for field %d: %w", id, err)
	}

	return nil
}

func (m *Message) dataFieldValue(id int) reflect.Value {
	return m.dataValue.FieldByName(fmt.Sprintf("F%d", id))
}
