package prefix

type Prefixer interface {
	// Returns field length encoded into []byte
	EncodeLength(maxLen, length int) ([]byte, error)

	// Returns the size of the field (number of characters, HEX-digits, bytes)
	// as well as the number of bytes read to decode the length
	DecodeLength(maxLen int, data []byte) (length int, read int, err error)

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
