package utils

const (
	// ErrInvalidEncoder is given when an invalid encoder is not supported
	ErrInvalidEncoder string = "invalid encoder"
	// ErrNonAvailableEncoding is given when encoding is non available
	ErrNonAvailableEncoding string = "non available encoding"
	// ErrInvalidLengthEncoder is given when the length of an encoder is invalid
	ErrInvalidLengthEncoder string = "invalid length encoder"
	// ErrInvalidLengthHead is given when the length of the head is invalid
	ErrInvalidLengthHead string = "invalid length head"
	// ErrValueTooLong is given when the length of the field is different from the length supplied
	ErrValueTooLong string = "length of value is longer than definition; type=%s, def_len=%d, len=%d"
	// ErrBadRaw is given when the raw data is malformed
	ErrBadRaw string = "bad raw data"
	// ErrBadElementData is given when the raw data is invalid data
	ErrBadElementData string = "bad element data"
	// ErrBadBinary is given when the length of the raw data is invalid
	ErrBadBinary string = "bad binary data"
	// ErrParseLengthFailed is given when the length of the binary data is invalid
	ErrParseLengthFailed string = "parse length head failed"
	// ErrNonExistSpecification is given when there is no specification
	ErrNonExistSpecification string = "don't exist specification"
	// ErrInvalidSpecification is given when there is invalid specification
	ErrInvalidSpecification string = "has invalid specification"
	// ErrInvalidBitmapArray is given when bitmap array is invalid
	ErrInvalidBitmapArray string = "invalid iso8583 bitmap array"
	// ErrInvalidElementLength is given when the length of the element is invalid
	ErrInvalidElementLength string = "invalid element length"
	// ErrInvalidElementType is given when the type of the element is invalid
	ErrInvalidElementType string = "invalid element type"
	// ErrMisMatchElementsBitmap is given when mismatch between bitmap and data elements
	ErrMisMatchElementsBitmap string = "don't match bitmap and data elements"
	// ErrNonInitializedMessage is given when message instance is not initialized
	ErrNonInitializedMessage string = "non initialized message"
)
