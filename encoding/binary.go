package encoding

import (
	"fmt"
)

var (
	_      Encoder = (*binaryEncoder)(nil)
	Binary         = &binaryEncoder{}
)

type binaryEncoder struct{}

func (e binaryEncoder) Encode(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)

	return out, nil
}

func (e binaryEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	if length > len(data) {
		return nil, 0, fmt.Errorf("failed to perform binary decoding: length %v exceeds the data size %v", length, len(data))
	}
	if length == 0 {
		return nil, 0, nil
	}

	out := make([]byte, length)
	copy(out, data[:length])

	return out, length, nil
}
