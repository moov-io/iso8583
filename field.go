package iso8583

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	ASCII = iota
	BCD
)

const (
	ERR_INVALID_ENCODER     string = "invalid encoder"
	ERR_INVALID_LENGTH_HEAD string = "invalid length head"
	ERR_MISSING_LENGTH      string = "missing length"
	ERR_VALUE_TOO_LONG      string = "length of value is longer than definition"
	ERR_BAD_RAW             string = "bad raw data"
	ERR_PARSE_LENGTH_FAILED string = "parse length head failed"
)

type Iso8583Type interface {
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
	if length == -1 {
		panic(ERR_MISSING_LENGTH)
	}
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
	if length == -1 {
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
	if length == -1 {
		panic(ERR_MISSING_LENGTH)
	}
	if len(val) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}
	if len(val) < length {
		val = append([]byte(strings.Repeat(" ", length-len(val))), val...)
	}
	return val, nil
}

func (a *Alphanumeric) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	if length == -1 {
		panic(ERR_MISSING_LENGTH)
	}
	a.Value = string(raw[:length])
	return length, nil
}

type Llvar struct {
	Value string
}

func NewLlvar(val string) *Llvar {
	return &Llvar{val}
}

func (l *Llvar) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val := []byte(l.Value)
	if length != -1 && len(val) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}
	if encoder != ASCII {
		panic(ERR_INVALID_ENCODER)
	}

	lenStr := fmt.Sprintf("%02d", len(val))
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch lenEncoder {
	case ASCII:
		lenVal = contentLen
		if len(lenVal) > 2 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
	case BCD:
		if len(lenVal) > 1 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
		lenVal = rbcd(contentLen)
	default:
		panic(ERR_INVALID_ENCODER)
	}
	return append(lenVal, val...), nil
}

func (l *Llvar) Load(raw []byte, encoder, lenEncoder, length int) (read int, err error) {
	// parse length head:
	var contentLen int
	switch lenEncoder {
	case ASCII:
		read = 2
		contentLen, err = strconv.Atoi(string(raw[:read]))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[:2]))
		}
	case BCD:
		read = 1
		contentLen, err = strconv.Atoi(string(bcdr2Ascii(raw[:read], 2)))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[0]))
		}
	default:
		panic(ERR_INVALID_ENCODER)
	}

	// parse body:
	body := raw[read : read+contentLen]
	read += contentLen
	if encoder != ASCII {
		panic(ERR_INVALID_ENCODER)
	}
	l.Value = string(body)

	return read, nil
}

type Llnumeric struct {
	Value string
}

func NewLlnumeric(val string) *Llnumeric {
	return &Llnumeric{val}
}

func (l *Llnumeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	raw := []byte(l.Value)
	if length != -1 && len(raw) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}

	val := raw
	switch encoder {
	case ASCII:
	case BCD:
		val = lbcd(raw)
	default:
		panic(ERR_INVALID_ENCODER)
	}

	lenStr := fmt.Sprintf("%02d", len(raw)) // length of digital characters
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch lenEncoder {
	case ASCII:
		lenVal = contentLen
		if len(lenVal) > 2 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
	case BCD:
		if len(lenVal) > 1 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
		lenVal = rbcd(contentLen)
	default:
		panic(ERR_INVALID_ENCODER)
	}
	return append(lenVal, val...), nil
}

func (l *Llnumeric) Load(raw []byte, encoder, lenEncoder, length int) (read int, err error) {
	// parse length head:
	var contentLen int
	switch lenEncoder {
	case ASCII:
		read = 2
		contentLen, err = strconv.Atoi(string(raw[:read]))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[:2]))
		}
	case BCD:
		read = 1
		contentLen, err = strconv.Atoi(string(bcdr2Ascii(raw[:read], 2)))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[0]))
		}
	default:
		panic(ERR_INVALID_ENCODER)
	}

	// parse body:
	switch encoder {
	case ASCII:
		l.Value = string(raw[read : read+contentLen])
		read += contentLen
	case BCD:
		bcdLen := (contentLen + 1) / 2
		l.Value = string(bcdl2Ascii(raw[read:read+bcdLen], contentLen))
		read += bcdLen
	default:
		panic(ERR_INVALID_ENCODER)
	}
	return read, nil
}

type Lllvar struct {
	Value string
}

func NewLllvar(val string) *Lllvar {
	return &Lllvar{val}
}

func (l *Lllvar) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val := []byte(l.Value)
	if length != -1 && len(val) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}

	switch encoder {
	case ASCII:
	case BCD:
		val = lbcd(val)
	default:
		panic(ERR_INVALID_ENCODER)
	}

	lenStr := fmt.Sprintf("%03d", len(val))
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch lenEncoder {
	case ASCII:
		lenVal = contentLen
		if len(lenVal) > 3 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
	case BCD:
		if len(lenVal) > 2 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
		lenVal = rbcd(contentLen)
	}
	return append(lenVal, val...), nil
}

func (l *Lllvar) Load(raw []byte, encoder, lenEncoder, length int) (read int, err error) {
	// parse length head:
	var contentLen int
	switch lenEncoder {
	case ASCII:
		read = 3
		contentLen, err = strconv.Atoi(string(raw[:read]))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[:3]))
		}
	case BCD:
		read = 2
		contentLen, err = strconv.Atoi(string(bcdr2Ascii(raw[:read], 2)))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[:2]))
		}
	default:
		panic(ERR_INVALID_ENCODER)
	}

	// parse body:
	body := raw[read : read+contentLen]
	read += contentLen
	if encoder != ASCII {
		panic(ERR_INVALID_ENCODER)
	}
	l.Value = string(body)

	return read, nil
}

type Lllnumeric struct {
	Value string
}

func NewLllnumeric(val string) *Lllnumeric {
	return &Lllnumeric{val}
}

func (l *Lllnumeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	raw := []byte(l.Value)
	if length != -1 && len(raw) > length {
		return nil, errors.New(ERR_VALUE_TOO_LONG)
	}

	val := raw
	switch encoder {
	case ASCII:
	case BCD:
		val = lbcd(raw)
	default:
		panic(ERR_INVALID_ENCODER)
	}

	lenStr := fmt.Sprintf("%03d", len(raw)) // length of digital characters
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch lenEncoder {
	case ASCII:
		lenVal = contentLen
		if len(lenVal) > 3 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
	case BCD:
		if len(lenVal) > 2 {
			panic(ERR_INVALID_LENGTH_HEAD)
		}
		lenVal = rbcd(contentLen)
	default:
		panic(ERR_INVALID_ENCODER)
	}
	return append(lenVal, val...), nil
}

func (l *Lllnumeric) Load(raw []byte, encoder, lenEncoder, length int) (read int, err error) {
	// parse length head:
	var contentLen int
	switch lenEncoder {
	case ASCII:
		read = 3
		contentLen, err = strconv.Atoi(string(raw[:read]))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[:3]))
		}
	case BCD:
		read = 2
		contentLen, err = strconv.Atoi(string(bcdr2Ascii(raw[:read], 2)))
		if err != nil {
			return 0, errors.New(ERR_PARSE_LENGTH_FAILED + ": " + string(raw[:2]))
		}
	default:
		panic(ERR_INVALID_ENCODER)
	}

	// parse body:
	switch encoder {
	case ASCII:
		l.Value = string(raw[read : read+contentLen])
		read += contentLen
	case BCD:
		bcdLen := (contentLen + 1) / 2
		l.Value = string(bcdl2Ascii(raw[read:read+bcdLen], contentLen))
		read += bcdLen
	default:
		panic(ERR_INVALID_ENCODER)
	}
	return read, nil
}
