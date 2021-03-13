package encoding

type Encoder interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte, int) ([]byte, error)
}
