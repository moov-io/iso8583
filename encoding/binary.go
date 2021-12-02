package encoding

import (
	"fmt"
)

var Binary Encoder = &binaryEncoder{}

type binaryEncoder struct{}

func (e binaryEncoder) Encode(data []byte) ([]byte, error) {
	out := append([]byte(nil), data...)

	return out, nil
}

func (e binaryEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	out := append([]byte(nil), data...)

	if length > len(data) {
		return nil, 0, fmt.Errorf("failed to perform binary decoding: length %v exceeds the data size %v", length, len(data))
	}

	return out[:length], length, nil
}
