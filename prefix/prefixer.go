package prefix

type Prefixer interface {
	// Returns field length encoded into []byte
	EncodeLength(maxLen, length int) ([]byte, error)

	// Retuns field length read from data
	DecodeLength(maxLen int, data []byte) (int, error)

	// Returns the number of bytes that takes to encode the length
	Length() int

	// Returns human readable infomation about length prefixer
	Inspect() string
}

type Prefixers struct {
	Fixed Prefixer
	L     Prefixer
	LL    Prefixer
	LLL   Prefixer
	LLLL  Prefixer
}

type PrefixerBuilder func(int) Prefixer
