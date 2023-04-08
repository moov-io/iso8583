package encoding

type Encoder interface {
	// Encode packs the data and returns the length in the encoding's representation.
	// the lenght is usually the number of bytes, but in bcd the number of encoded digits.
	Encode([]byte) (packed []byte, length int, err error)
	// Decode returns data decoded into ASCII (or bytes), how many bytes were read, error
	Decode([]byte, int) (data []byte, read int, err error)
}
