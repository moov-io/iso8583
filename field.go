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
	ERR_INVALID_ENCODER        string = "invalid encoder"
	ERR_INVALID_LENGTH_ENCODER string = "invalid length encoder"
	ERR_INVALID_LENGTH_HEAD    string = "invalid length head"
	ERR_MISSING_LENGTH         string = "missing length"
	ERR_VALUE_TOO_LONG         string = "length of value is longer than definition; type=%s, def_len=%d, len=%d"
	ERR_BAD_RAW                string = "bad raw data"
	ERR_PARSE_LENGTH_FAILED    string = "parse length head failed"
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
