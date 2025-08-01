package prefix

import "errors"

const fieldLengthIsLargerThanMax = "data length: %d is larger than maximum: %d"
const numberOfDigitsInLengthExceeds = "number of digits in length: %d exceeds: %d"
const notEnoughDataToRead = "not enough data length: %d to read: %d byte digits"
const invalidLength = "invalid length: %d"
const dataLengthIsLargerThanMax = "data length: %d is larger than maximum %d"
const fieldLengthShouldBeFixed = "data length: %d should be fixed: %d"

type LengthError struct {
	err error
}

func (e *LengthError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func IsLengthError(err error) bool {
	e := &LengthError{}
	if errors.As(err, &e) {
		return true
	}
	return false
}
