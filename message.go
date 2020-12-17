package iso8583

import (
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/spec"
	"github.com/moov-io/iso8583/utils"
)

type Message struct {
	Fields map[int]Field
	spec   *spec.MessageSpec

	// let's keep it 8 bytes for now
	bitmap *utils.Bitmap
	Data   interface{}
}

type Setter interface {
	Set(b []byte)
}

func NewMessage(spec *spec.MessageSpec) *Message {
	return &Message{
		Fields: map[int]Field{},
		spec:   spec,
	}
}

func (m *Message) Set(id int, field Field) {
	m.Fields[id] = field
}

func (m *Message) Bitmap() *utils.Bitmap {
	return m.bitmap
}

func (m *Message) MTI(val string) {
	m.Fields[0] = NewField(0, []byte(val))
}

func (m *Message) Field(id int, val string) {
	m.Fields[id] = NewField(id, []byte(val))
}

func (m *Message) BinaryField(id int, val []byte) {
	m.Fields[id] = NewField(id, val)
}

func (m *Message) GetMTI() string {
	// check index
	return m.Fields[0].String()
}

func (m *Message) GetString(id int) string {
	if field, ok := m.Fields[id]; ok {
		return field.String()
	}

	return ""
}

func (m *Message) GetBytes(id int) []byte {
	if field, ok := m.Fields[id]; ok {
		return field.Bytes()
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
	packedMTI, err := m.spec.Fields[0].Pack(m.Fields[0].Bytes())
	if err != nil {
		return nil, err
	}
	packed = append(packed, packedMTI...)

	// pack Bitmap
	packedBitmap, err := m.spec.Fields[1].Pack(m.bitmap.Bytes())
	if err != nil {
		return nil, err
	}
	packed = append(packed, packedBitmap...)

	// pack each field
	for i := 2; i <= maxId; i++ {
		if field, ok := m.Fields[i]; ok {
			def, ok := m.spec.Fields[i]
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
	data, read, err := m.spec.Fields[0].Unpack(src)
	if err != nil {
		return err
	}
	m.BinaryField(0, data)

	off = read

	data, read, err = m.spec.Fields[1].Unpack(src[off:])
	if err != nil {
		return err
	}
	m.BinaryField(1, data)
	m.bitmap = utils.NewBitmapFromData(data)
	off += read

	for i := 2; i <= m.bitmap.Len(); i++ {
		if m.bitmap.IsSet(i) {
			fieldSpec, ok := m.spec.Fields[i]
			if !ok {
				return fmt.Errorf("Failed to unpack field %d. No Specification found for the field", i)
			}

			data, read, err = fieldSpec.Unpack(src[off:])
			if err != nil {
				return err
			}

			m.setDataFieldValue(i, data)
			m.BinaryField(i, data)
			off += read
		}
	}

	return nil
}

func (m *Message) setDataFieldValue(i int, data []byte) {
	if m.Data == nil {
		return
	}

	// get the struct
	str := reflect.ValueOf(m.Data).Elem()

	// get the struct field
	field := str.FieldByName(fmt.Sprintf("F%d", i))
	if field == (reflect.Value{}) {
		return
	}

	// we should check if it's nil
	// if it's not nil, let's create a field
	// with a corresponding type

	// get the type of the field
	fieldType := field.Type().Elem()

	// create new field
	fieldVal := reflect.New(fieldType)
	field.Addr().Elem().Set(fieldVal)

	// setting value may happen after we add created field to the struct
	// as we work with pointers

	// get field as a Setter
	setterField := fieldVal.Interface().(Setter)

	// finally set value to the field
	setterField.Set(data)
}
