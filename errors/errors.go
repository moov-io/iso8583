package errors

import (
	"errors"
)

// UnpackError returns an error with the possibility to access the RawMessage when
// the connection failed to unpack the message
type UnpackError struct {
	Err        error
	FieldID    string
	RawMessage []byte
}

func (e *UnpackError) Error() string {
	return e.Err.Error()
}

func (e *UnpackError) Unwrap() error {
	return e.Err
}

// FieldIDs returns the list of field and subfield IDs (if any) that errored from outermost inwards
func (e *UnpackError) FieldIDs() []string {
	fieldIDs := []string{e.FieldID}
	err := e.Err
	var unpackError *UnpackError
	for {
		if errors.As(err, &unpackError) {
			fieldIDs = append(fieldIDs, unpackError.FieldID)
			err = unpackError.Unwrap()
		} else {
			break
		}
	}

	return fieldIDs
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
