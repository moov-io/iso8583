package iso8583

import (
	"errors"
	"strings"
)

const (
	ASCII = iota
	BCD
)

const (
	ERR_INVALID_ENCODER string = "invalid encoder"
	ERR_MISSING_LENGTH = "missing length"
	ERR_VALUE_TOO_LONG string = "length of value is longer than definition"
	ERR_BAD_RAW = "bad raw data"
)

type Iso8583Type interface{
	// Byte representation of current field.
	Bytes(encoder, lenEncoder, length int) ([]byte, error)
	
	// Load unmarshal byte value into Iso8583Type according to the
	// specific arguments. It returns the number of bytes actually read. 
	Load(raw []byte, encoder, lenEncoder, length int) (int, error)
}

// A Numeric contains numeric value only in fix length. It holds numeric
// value as a string. Supportted encoder are ascii and bcd. Length is
// required for marshalling and unmarshalling.
type Numeric struct {
	Value string
}

func NewNumeric(val string) *Numeric {
	return &Numeric{val}
}

func (n *Numeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val := []byte(n.Value)
	if len(val) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}
	if len(val) < length {
		val = append([]byte(strings.Repeat("0", length-len(val))), val...)
	}
	switch encoder {
	case BCD:
		return lbcd(val), nil
	case ASCII:
		return val, nil
	default:
		panic(ERR_INVALID_ENCODER)
	}
}

func (n *Numeric) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	if length < 1 {
		panic(ERR_MISSING_LENGTH)
	}
	switch encoder {
	case BCD:
		l := (length + 1) / 2
		if len(raw) < l {
			return 0, errors.New(ERR_BAD_RAW)
		}
		n.Value = string(bcdl2Ascii(raw[:l], length))
		return l, nil
	case ASCII:
		if len(raw) < length {
			return 0, errors.New(ERR_BAD_RAW)
		}
		n.Value = string(raw[:length])
		return length, nil
	default:
		panic(ERR_INVALID_ENCODER)
	}
}

// An Alphanumeric contains alphanumeric value in fix length. The only
// supportted encoder is ascii. Length is required for marshalling and
// unmarshalling.
type Alphanumeric struct {
 	Value string
}

 func NewAlphanumeric(val string) *Alphanumeric {
 	return &Alphanumeric{Value: val}
 }

func (a *Alphanumeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val := []byte(a.Value)
	if len(val) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}
	if len(val) < length {
		val = append([]byte(strings.Repeat(" ", length-len(val))), val...)
	}
	return val, nil
}

func (a *Alphanumeric) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	if length < 1 {
		panic(ERR_MISSING_LENGTH)
	}
	a.Value = string(raw[:length])
	return length, nil
}

type Llvar struct {
	Value string
}

func (l *Llvar) Bytes(encoder, lenEncoder, length int) ([]byte, error) {

	return nil, nil
}

func (l *Llvar) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	return length, nil
}
	
type Lllvar struct {
	Value string
}

func (l *Lllvar) Bytes(encoder, lenEncoder, length int) ([]byte, error) {

	return nil, nil
}

func (l *Lllvar) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	return length, nil
}

