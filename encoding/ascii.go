package encoding

import (
	"fmt"
	"io"
)

var ASCII Coder = &asciiEncoder{}

type asciiEncoder struct{}

func (e asciiEncoder) Encode(data []byte) ([]byte, error) {
	out := []byte{}
	for _, r := range data {
		if r > 127 {
			return nil, fmt.Errorf("invalid ASCII char: '%s'", string(r))
		}
		out = append(out, r)
	}

	return out, nil
}

func (e asciiEncoder) DecodeFrom(r io.Reader, length int) (data []byte, read int, err error) {
	data = make([]byte, length)
	read, err = io.ReadFull(r, data)
	if err != nil {
		return nil, read, fmt.Errorf("reading %d bytes from reader: %v", length, err)
	}

	for _, r := range data {
		if r > 127 {
			return nil, 0, fmt.Errorf("invalid ASCII char: '%s'", string(r))
		}
	}

	return data, read, nil
}

func (e asciiEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	// read only 'length' bytes (1 byte - 1 ASCII character)
	data = data[:length]
	out := []byte{}
	for _, r := range data {
		if r > 127 {
			return nil, 0, fmt.Errorf("invalid ASCII char: '%s'", string(r))
		}
		out = append(out, r)
	}

	return out, length, nil
}
