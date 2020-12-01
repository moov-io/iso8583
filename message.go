package iso8583

import (
	"fmt"

	"github.com/moov-io/iso8583/fields"
	"github.com/moov-io/iso8583/spec"
	"github.com/moov-io/iso8583/utils"
)

type Message struct {
	Fields map[int]fields.Field
	spec   *spec.MessageSpec

	// let's keep it 8 bytes for now
	bitmap *utils.Bitmap
}

func NewMessage(spec *spec.MessageSpec) *Message {
	return &Message{
		Fields: map[int]fields.Field{},
		spec:   spec,
	}
}

func (m *Message) Set(id int, field fields.Field) {
	m.Fields[id] = field
}

func (m *Message) Field(id int, val string) {
	m.Fields[id] = fields.NewField(id, []byte(val))
}

func (m *Message) BinaryField(id int, val []byte) {
	m.Fields[id] = fields.NewField(id, val)
}

func (m *Message) GetString(id int) string {
	// check index
	return m.Fields[id].String()
}

func (m *Message) GetBytes(id int) []byte {
	// check index
	return m.Fields[id].Bytes()
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

	// unpack Bitmap
	fmt.Println("offset: ", off)
	data, read, err = m.spec.Fields[1].Unpack(src[off:])
	if err != nil {
		return err
	}
	m.BinaryField(1, data)
	m.bitmap = utils.NewBitmapFromData(data)
	off += read

	// we should take data from the bitmap
	// but for simplicity

	fmt.Println("offset: ", off)
	fmt.Println("unpacked bitmap:", data)

	// read MTI
	// unpack Bitmap
	// based on bitmap length unpack all fields

	for i := 2; i <= m.bitmap.Len(); i++ {
		fmt.Printf("Checking field: %d and it is: %v\n", i, m.bitmap.IsSet(i))
		if m.bitmap.IsSet(i) {
			fmt.Printf("Data at offset %d: %v\n", off, src[off:])
			data, read, err = m.spec.Fields[i].Unpack(src[off:])
			if err != nil {
				return err
			}
			fmt.Printf("Read data for field: %d =>%v\n", i, data)
			m.BinaryField(i, data)
			off += read
		}
	}

	return nil
}
