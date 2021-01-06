package field

type Field interface {
	Spec() *Spec
	SetSpec(spec *Spec)
	Pack(data []byte) ([]byte, error)
	Unpack(data []byte) ([]byte, int, error)

	SetBytes(b []byte)
	Bytes() []byte

	String() string
}
