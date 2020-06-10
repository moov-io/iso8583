package iso8583

import (
	"errors"
	"fmt"
	"strings"
)

// An Alphanumeric contains alphanumeric value in fix length. The only
// supportted encoder is ascii. Length is required for marshalling and
// unmarshalling.
type Alphanumeric struct {
	Value string
}

// NewAlphanumeric create new Alphanumeric field
func NewAlphanumeric(val string) *Alphanumeric {
	return &Alphanumeric{Value: val}
}

// IsEmpty check Alphanumeric field for empty value
func (a *Alphanumeric) IsEmpty() bool {
	return len(a.Value) == 0
}

// Bytes encode Alphanumeric field to bytes
func (a *Alphanumeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val := []byte(a.Value)
	val, err := UTF8ToWindows1252(val)
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, errors.New(ERR_MISSING_LENGTH)
	}
	if len(val) > length {
		return nil, fmt.Errorf(ERR_VALUE_TOO_LONG, "Alphanumeric", length, len(val))
	}
	if len(val) < length {
		val = append([]byte(strings.Repeat(" ", length-len(val))), val...)
	}
	return val, nil
}

// Load decode Alphanumeric field from bytes
func (a *Alphanumeric) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	raw, err := UTF8ToWindows1252(raw)
	if err != nil {
		return 0, err
	}
	if length == -1 {
		return 0, errors.New(ERR_MISSING_LENGTH)
	}
	if len(raw) < length {
		return 0, errors.New(ERR_BAD_RAW)
	}
	a.Value = string(raw[:length])
	return length, nil
}
