package prefixer

type Prefixer interface {
	// Returns field length encoded into []byte
	EncodeLength(length int) ([]byte, error)

	// Retuns field length read from data
	DecodeLength(data []byte) (int, error)

	// Returns the number of bytes that takes to encode the length
	Length() int

	// Returns human readable infomation about length prefixer
	Inspect() string
}

type Prefixers struct {
	Fixed PrefixerBuilder
	L     PrefixerBuilder
	LL    PrefixerBuilder
	LLL   PrefixerBuilder
	LLLL  PrefixerBuilder
}

type PrefixerBuilder func(int) Prefixer
