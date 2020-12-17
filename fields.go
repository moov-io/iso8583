package iso8583

import "strconv"

type Field interface {
	String() string
	Bytes() []byte
}

type binaryField struct {
	ID  int
	val []byte
}

func (f *binaryField) String() string {
	return string(f.val)
}

func (f *binaryField) Bytes() []byte {
	return f.val
}

func NewField(id int, val []byte) Field {
	return &binaryField{
		ID:  id,
		val: val,
	}
}

type StringField struct {
	Value string
}

func (f *StringField) Set(b []byte) {
	f.Value = string(b)
}

type NumericField struct {
	Value int
}

func (f *NumericField) Set(b []byte) {
	f.Value, _ = strconv.Atoi(string(b))
}
