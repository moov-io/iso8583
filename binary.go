package iso8583

import (
	"errors"
	"fmt"
)

// Binary contains binary value
type Binary struct {
	Value  []byte
	FixLen int
}

// NewBinary create new Binary field
func NewBinary(d []byte) *Binary {
	return &Binary{d, -1}
}

// IsEmpty check Binary field for empty value
func (b *Binary) IsEmpty() bool {
	return len(b.Value) == 0
}

// Bytes encode Binary field to bytes
func (b *Binary) Bytes(encoder, lenEncoder, l int) ([]byte, error) {
	length := l
	if b.FixLen != -1 {
		length = b.FixLen
	}
	if length == -1 {
		return nil, errors.New(ErrMissingLength)
	}
	if len(b.Value) > length {
		return nil, fmt.Errorf(ErrValueTooLong, "Binary", length, len(b.Value))
	}
	if len(b.Value) < length {
		return append(b.Value, make([]byte, length-len(b.Value))...), nil
	}
	return b.Value, nil
}

// Load decode Binary field from bytes
func (b *Binary) Load(raw []byte, encoder, lenEncoder, length int) (int, error) {
	if length == -1 {
		return 0, errors.New(ErrMissingLength)
	}
	if len(raw) < length {
		return 0, errors.New(ErrBadRaw)
	}
	b.Value = raw[:length]
	b.FixLen = length
	return length, nil
}
