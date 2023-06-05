package encoding

var (
	AsciiTLVTagL    = &asciiTLVTagL{}
	AsciiTLVTagLL   = &asciiTLVTagLL{}
	AsciiTLVTagLLL  = &asciiTLVTagLLL{}
	AsciiTLVTagLLLL = &asciiTLVTagLLLL{}
)

type (
	asciiTLVTagL    struct{ asciiEncoder }
	asciiTLVTagLL   struct{ asciiEncoder }
	asciiTLVTagLLL  struct{ asciiEncoder }
	asciiTLVTagLLLL struct{ asciiEncoder }
)

func (e asciiTLVTagL) Decode(data []byte, _ int) ([]byte, int, error) {
	return e.asciiEncoder.Decode(data, 1)
}

func (e asciiTLVTagLL) Decode(data []byte, _ int) ([]byte, int, error) {
	return e.asciiEncoder.Decode(data, 2)
}

func (e asciiTLVTagLLL) Decode(data []byte, _ int) ([]byte, int, error) {
	return e.asciiEncoder.Decode(data, 3)
}

func (e asciiTLVTagLLLL) Decode(data []byte, _ int) ([]byte, int, error) {
	return e.asciiEncoder.Decode(data, 4)
}
