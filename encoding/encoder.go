package encoding

type Encoder interface {
	// Encode encodes source data (ASCII, characters, digits, etc.) into
	// destination bytes. It returns encoded bytes and any error
	Encode([]byte) ([]byte, error)

	// Decode decodes data into into bytes (ASCII, characters, digits,
	// etc.). It returns the bytes representing the decoded data, the
	// number of bytes read from the input, and any error
	Decode([]byte, int) (data []byte, read int, err error)
}
