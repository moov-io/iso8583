package encoding

import (
	"bytes"
	"io"
)

var Binary Coder = &binaryEncoder{}

type binaryEncoder struct{}

func (e binaryEncoder) Encode(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)

	return out, nil
}

func (e binaryEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	return e.DecodeFrom(bytes.NewReader(data), length)
}

func (e binaryEncoder) DecodeFrom(r io.Reader, length int) (data []byte, read int, err error) {
	data = make([]byte, length)
	read, err = io.ReadFull(r, data)

	return data, read, err
}
