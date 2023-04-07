package encoding

import (
	"fmt"
)

var (
	_      Encoder = (*binaryEncoder)(nil)
	Binary         = &binaryEncoder{}
)

type binaryEncoder struct{}

func (e binaryEncoder) Encode(data []byte) ([]byte, int, error) {
	out := append([]byte(nil), data...)

	return out, len(out), nil
}

func (e binaryEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	out := append([]byte(nil), data...)

	if length > len(data) {
		return nil, 0, fmt.Errorf("failed to perform binary decoding: length %v exceeds the data size %v", length, len(data))
	}

	return out[:length], length, nil
}
