package encoding

var (
	AsciiTLVTagL    = &asciiTLVTag{digits: 1}
	AsciiTLVTagLL   = &asciiTLVTag{digits: 2}
	AsciiTLVTagLLL  = &asciiTLVTag{digits: 3}
	AsciiTLVTagLLLL = &asciiTLVTag{digits: 4}
)

type asciiTLVTag struct {
	asciiEncoder
	digits int
}

func (e asciiTLVTag) Decode(data []byte, _ int) ([]byte, int, error) {
	return e.asciiEncoder.Decode(data, e.digits)
}
