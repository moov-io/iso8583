package field

type Field interface {
	Spec() *Spec
	SetSpec(spec *Spec)

	Pack() ([]byte, error)
	Unpack(data []byte) (int, error)

	SetBytes(b []byte) error
	Bytes() ([]byte, error)

	SetData(d interface{}) error

	String() (string, error)
}
