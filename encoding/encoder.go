package encoding

type Encoder interface {
	Encode([]byte) (packed []byte, length int, err error)
	// Returns data decoded into ASCII (or bytes), how many bytes were read, error
	Decode([]byte, int) (data []byte, read int, err error)
}
