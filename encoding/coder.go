package encoding

import "io"

type Coder interface {
	// Encode encodes src data into encoding, returns encoded data and error
	Encode(src []byte) (dst []byte, err error)

	// Decode decodes src data of the length specified for encoding (ASCII,
	// BCD, HEX) into ASCII. Returns decoded date, how many bytes were
	// read and error
	Decode(src []byte, length int) (data []byte, read int, err error)

	// Reads data of the length specified for encoding (ASCII, BCD, HEX)
	// from Reader and decodes it into ASCII (or bytes). Returns decoded
	// date, how many bytes were read, error
	// examples of length:
	// - 2 digits BCD - 1 byte to read, returns 2 bytes in ASCII
	// - 4 digits ASCII - 4 bytes to read, returns 3 bytes in ASCII
	DecodeFrom(r io.Reader, length int) (data []byte, read int, err error)
}
