package prefix

type Prefixer interface {
	// Returns field length encoded into []byte
	EncodeLength(maxLen, length int) ([]byte, error)

	// Returns the size of the field (number of characters, HEX-digits, bytes)
	DecodeLength(maxLen int, data []byte) (int, error)

	// Returns the number of bytes that takes to encode the length
	Length() int

	// Returns human readable information about length prefixer. Returned value
	// is used to create prefixer when we build spec from a JSON spec.
	// Returned value should be in the following format:
	//  PrefixerName.Length
	// Examples:
	//  ASCII.LL
	//  Hex.Fixed
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
