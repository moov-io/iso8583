package encoding

type Encoder interface {
	Encode([]byte) ([]byte, error)
	// Returns data decoded into ASCII (or bytes), how many bytes were read, error
	Decode([]byte, int) (data []byte, read int, err error)
}
