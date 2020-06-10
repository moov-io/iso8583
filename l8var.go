// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package iso8583

import (
	"errors"
	"fmt"
	"strconv"
)

// L8var contains bytes in non-fixed length field, first 8 symbols of field contains length
type L8var struct {
	Value []byte
}

// NewL8var create new L8var field
func NewL8var(val []byte) *L8var {
	return &L8var{val}
}

// IsEmpty check Lllvar field for empty value
func (l *L8var) IsEmpty() bool {
	return len(l.Value) == 0
}

// Bytes encode Lllvar field to bytes
func (l *L8var) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	val, err := UTF8ToWindows1252(l.Value)
	if err != nil {
		return nil, err
	}

	if length != -1 && len(val) > length {
		return nil, fmt.Errorf(ErrValueTooLong, "L8var", length, len(val))
	}
	if encoder != ASCII {
		return nil, errors.New(ErrInvalidEncoder)
	}

	lenStr := fmt.Sprintf("%08d", len(val))
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch lenEncoder {
	case ASCII:
		lenVal = contentLen
		if len(lenVal) > 8 {
			return nil, errors.New(ErrInvalidLengthHead)
		}
	case rBCD:
		fallthrough
	case BCD:
		lenVal = rbcd(contentLen)
		if len(lenVal) > 7 || len(contentLen) > 8 {
			return nil, errors.New(ErrInvalidLengthHead)
		}
	default:
		return nil, errors.New(ErrInvalidLengthEncoder)
	}
	return append(lenVal, val...), nil
}

// Load decode Lllvar field from bytes
func (l *L8var) Load(raw []byte, encoder, lenEncoder, length int) (read int, err error) {
	raw, err = UTF8ToWindows1252(raw)
	if err != nil {
		return 0, err
	}

	// parse length head:
	var contentLen int
	switch lenEncoder {
	case ASCII:
		read = 8
		contentLen, err = strconv.Atoi(string(raw[:read]))
		if err != nil {
			return 0, errors.New(ErrParseLengthFailed + ": " + string(raw[:8]))
		}
	case rBCD:
		fallthrough
	case BCD:
		read = 7
		contentLen, err = strconv.Atoi(string(bcdr2Ascii(raw[:read], 8)))
		if err != nil {
			return 0, errors.New(ErrParseLengthFailed + ": " + string(raw[:8]))
		}
	default:
		return 0, errors.New(ErrInvalidLengthEncoder)
	}
	if len(raw) < (read + contentLen) {
		return 0, errors.New(ErrBadRaw)
	}
	// parse body:
	l.Value = raw[read : read+contentLen]
	read += contentLen
	if encoder != ASCII {
		return 0, errors.New(ErrInvalidEncoder)
	}

	return read, nil
}
