package field

import "io"

type Field interface {
	Spec() *Spec
	SetSpec(spec *Spec)

	Pack() ([]byte, error)

	ReadFrom(io.Reader) (int, error)

	SetBytes(b []byte) error
	Bytes() ([]byte, error)

	SetData(d interface{}) error

	String() (string, error)
}
