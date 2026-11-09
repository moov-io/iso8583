package encoding

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
)

var (
	_     Encoder = (*asciiEncoder)(nil)
	ASCII         = &asciiEncoder{}
)

type asciiEncoder struct{}

func (e asciiEncoder) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	out := make([]byte, len(data))
	for i, r := range data {
		if r > 127 {
			return nil, utils.NewSafeError(fmt.Errorf("invalid ASCII char: '%s'", string(r)), "failed to perform ASCII encoding")
		}
		out[i] = r
	}

	return out, nil
}

func (e asciiEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	// length should be positive
	if length < 0 {
		return nil, 0, fmt.Errorf("invalid length: %d", length)
	}

	// read only 'length' bytes (1 byte - 1 ASCII character)
	if len(data) < length {
		return nil, 0, fmt.Errorf("not enough data to decode. expected len %d, got %d", length, len(data))
	}
	if length == 0 {
		return nil, 0, nil
	}

	data = data[:length]
	out := make([]byte, length)
	for i, r := range data {
		if r > 127 {
			return nil, 0, utils.NewSafeError(fmt.Errorf("invalid ASCII char: '%s'", string(r)), "failed to perform ASCII decoding")
		}
		out[i] = r
	}

	return out, length, nil
}
