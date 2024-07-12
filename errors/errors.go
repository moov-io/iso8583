package errors

// UnpackError returns an error with the possibility to access the RawMessage when
// the connection failed to unpack the message
type UnpackError struct {
	Err error
	// the field ID of the field on which unpacking errored
	FieldID string
	// the field ID and subfield IDs (if any) that errored ordered from outermost inwards
	FieldIDs   []string
	RawMessage []byte
}

func (e *UnpackError) Error() string {
	return e.Err.Error()
}

func (e *UnpackError) Unwrap() error {
	return e.Err
}

type PackError struct {
	Err error
}

func (e *PackError) Error() string {
	return e.Err.Error()
}

func (e *PackError) Unwrap() error {
	return e.Err
}
