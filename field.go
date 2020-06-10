package iso8583

const (
	// ASCII is ASCII encoding
	ASCII = iota
	// BCD is "left-aligned" BCD
	BCD
	// rBCD is "right-aligned" BCD with odd length (for ex. "643" as [6 67] == "0643"), only for Numeric, Llnumeric and Lllnumeric fields
	rBCD
)

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

// Iso8583Type interface for ISO 8583 fields
type Iso8583Type interface {
	// Byte representation of current field.
	Bytes(encoder, lenEncoder, length int) ([]byte, error)

	// Load unmarshal byte value into Iso8583Type according to the
	// specific arguments. It returns the number of bytes actually read.
	Load(raw []byte, encoder, lenEncoder, length int) (int, error)

	// IsEmpty check is field empty
	IsEmpty() bool
}
