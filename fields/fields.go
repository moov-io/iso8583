package fields

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
