package iso8583

import (
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/utils"
)

type Message struct {
	// should we make it private?
	Fields map[int]field.Field
	spec   *MessageSpec

	// let's keep it 8 bytes for now
	bitmap *utils.Bitmap
	Data   interface{}

	fieldsMap map[int]struct{}
}

func NewMessage(spec *MessageSpec) *Message {
	fields := spec.CreateMessageFields()

	return &Message{
		Fields:    fields,
		spec:      spec,
		fieldsMap: map[int]struct{}{},
	}
}

func (m *Message) Bitmap() *utils.Bitmap {
	return m.bitmap
}

func (m *Message) MTI(val string) {
	m.fieldsMap[0] = struct{}{}
	m.Fields[0].SetBytes([]byte(val))
}

func (m *Message) Field(id int, val string) {
	m.fieldsMap[id] = struct{}{}
	m.Fields[id].SetBytes([]byte(val))
}

func (m *Message) BinaryField(id int, val []byte) {
	m.fieldsMap[id] = struct{}{}
	m.Fields[id].SetBytes(val)
}

func (m *Message) GetMTI() string {
	// check index
	return m.Fields[0].String()
}

func (m *Message) GetString(id int) string {
	if _, ok := m.fieldsMap[id]; ok {
		return m.Fields[id].String()
	}

	return ""
}

func (m *Message) GetBytes(id int) []byte {
	if _, ok := m.fieldsMap[id]; ok {
		return m.Fields[id].Bytes()
	}

	return nil
}

func (m *Message) Pack() ([]byte, error) {
	packed := []byte{}

	// use fixed length of the bitmap for now
	m.bitmap = utils.NewBitmap(128)

	// fill in the bitmap
	// and find max field id (to determine bitmap length later)
	maxId := 0
	for id, _ := range m.Fields {
		if id > maxId {
			maxId = id
		}

		// regular fields start from index 2
		if id < 2 {
			continue
		}
		m.bitmap.Set(id)
	}

	// pack MTI
	packedMTI, err := m.Fields[0].Pack(m.Fields[0].Bytes())
	if err != nil {
		return nil, err
	}
	packed = append(packed, packedMTI...)

	// pack Bitmap
	packedBitmap, err := m.Fields[1].Pack(m.bitmap.Bytes())
	if err != nil {
		return nil, err
	}
	packed = append(packed, packedBitmap...)

	// pack each field
	for i := 2; i <= maxId; i++ {
		if field, ok := m.Fields[i]; ok {
			def, ok := m.Fields[i]
			if !ok {
				return nil, fmt.Errorf("Failed to pack field: %d. No definition found", i)
			}

			packedField, err := def.Pack(field.Bytes())
			if err != nil {
				return nil, err
			}
			packed = append(packed, packedField...)
		}
	}

	return packed, nil
}

func (m *Message) Unpack(src []byte) error {
	var off int

	// unpack MTI
	data, read, err := m.Fields[0].Unpack(src)
	if err != nil {
		return err
	}
	m.BinaryField(0, data)

	off = read

	// hm... how to tell that this field was set?
	m.fieldsMap[1] = struct{}{}
	data, read, err = m.Fields[1].Unpack(src[off:])
	if err != nil {
		return err
	}
	m.bitmap = utils.NewBitmapFromData(data)
	off += read

	for i := 2; i <= m.bitmap.Len(); i++ {
		if m.bitmap.IsSet(i) {
			fl, ok := m.Fields[i]
			if !ok {
				return fmt.Errorf("Failed to unpack field %d. No Specification found for the field", i)
			}

			m.fieldsMap[i] = struct{}{}
			_, read, err = fl.Unpack(src[off:])
			if err != nil {
				return err
			}

			err = m.linkDataFieldWithMessageField(i, fl)
			if err != nil {
				return err
			}
			// m.BinaryField(i, data)
			off += read
		}
	}

	return nil
}

func (m *Message) linkDataFieldWithMessageField(i int, fl field.Field) error {
	if m.Data == nil {
		return nil
	}

	// get the struct
	str := reflect.ValueOf(m.Data).Elem()

	fieldName := fmt.Sprintf("F%d", i)

	// get the struct field
	dataField := str.FieldByName(fieldName)
	if dataField == (reflect.Value{}) {
		return nil
	}

	if dataField.Type() != reflect.TypeOf(fl) {
		return fmt.Errorf("field %s type: %v does not match the type in the spec: %v", fieldName, dataField.Type(), reflect.TypeOf(fl))
	}

	dataField.Addr().Elem().Set(reflect.ValueOf(fl))

	return nil
}
