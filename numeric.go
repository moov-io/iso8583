package iso8583

import (
	"errors"
	"fmt"
	"strings"
)

// A Numeric contains numeric value only in fix length. It holds numeric
// value as a string. Supportted encoder are ascii, bcd and rbcd. Length is
// required for marshalling and unmarshalling.
type Numeric struct {
	Value string
}

// NewNumeric create new Numeric field
func NewNumeric(val string) *Numeric {
	return &Numeric{val}
}

// IsEmpty check Numeric field for empty value
func (n *Numeric) IsEmpty() bool {
	return len(n.Value) == 0
}

// Bytes encode Numeric field to bytes
func (n *Numeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val := []byte(n.Value)
	if length == -1 {
		return nil, errors.New(ERR_MISSING_LENGTH)
	}
	// if encoder == rBCD then length can be, for example, 3,
	// but value can be, for example, "0631" (after decode from rBCD, because BCD use 1 byte for 2 digits),
	// and we can encode it only if first digit == 0
	if (encoder == rBCD) &&
		len(val) == (length+1) &&
		(string(val[0:1]) == "0") {
		// Cut value to length
		val = val[1:len(val)]
	}

	if len(val) > length {
		return nil, fmt.Errorf(ERR_VALUE_TOO_LONG, "Numeric", length, len(val))
	}
	if len(val) < length {
		val = append([]byte(strings.Repeat("0", length-len(val))), val...)
	}
	switch encoder {
	case BCD:
		return lbcd(val), nil
	case rBCD:
		return rbcd(val), nil
	case ASCII:
		return val, nil
	default:
		return nil, errors.New(ERR_INVALID_ENCODER)
	}
}

// Load decode Numeric field from bytes
func (n *Numeric) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	if length == -1 {
		return 0, errors.New(ERR_MISSING_LENGTH)
	}
	switch encoder {
	case BCD:
		l := (length + 1) / 2
		if len(raw) < l {
			return 0, errors.New(ERR_BAD_RAW)
		}
		n.Value = string(bcdl2Ascii(raw[:l], length))
		return l, nil
	case rBCD:
		l := (length + 1) / 2
		if len(raw) < l {
			return 0, errors.New(ERR_BAD_RAW)
		}
		n.Value = string(bcdr2Ascii(raw[0:l], length))
		return l, nil
	case ASCII:
		if len(raw) < length {
			return 0, errors.New(ERR_BAD_RAW)
		}
		n.Value = string(raw[:length])
		return length, nil
	default:
		return 0, errors.New(ERR_INVALID_ENCODER)
	}
}
