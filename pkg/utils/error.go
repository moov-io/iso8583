package utils

const (
	// ErrInvalidEncoder is given when an invalid encoder is not supported
	ErrInvalidEncoder string = "invalid encoder"
	// ErrInvalidLengthEncoder is given when the length of an encoder is invalid
	ErrInvalidLengthEncoder string = "invalid length encoder"
	// ErrInvalidLengthHead is given when the length of the head is invalid
	ErrInvalidLengthHead string = "invalid length head"
	// ErrMissingLength is given when the length of a field is missing
	ErrMissingLength string = "missing length"
	// ErrValueTooLong is given when the length of the field is different from the length supplied
	ErrValueTooLong string = "length of value is longer than definition; type=%s, def_len=%d, len=%d"
	// ErrBadRaw is given when the raw data is malformed
	ErrBadRaw string = "bad raw data"
	// ErrParseLengthFailed is given when the length of the raw data is invalid
	ErrParseLengthFailed string = "parse length head failed"
)
