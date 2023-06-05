package encoding

var (
	ASCIITLVTagL    = &asciiTLVTagL{}
	ASCIITLVTagLL   = &asciiTLVTagLL{}
	ASCIITLVTagLLL  = &asciiTLVTagLLL{}
	ASCIITLVTagLLLL = &asciiTLVTagLLLL{}
)

type (
	asciiTLVTagL    struct{ asciiEncoder }
	asciiTLVTagLL   struct{ asciiEncoder }
	asciiTLVTagLLL  struct{ asciiEncoder }
	asciiTLVTagLLLL struct{ asciiEncoder }
)
