package iso8583

// UnpackError returns error with possibility to access RawMessage when
// connection failed to unpack message
type UnpackError struct {
	Err        error
	Field      string
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
