package encoding

import (
	"encoding/hex"
)

var Hex Encoder = &hexEncoder{}

type hexEncoder struct{}

func (e hexEncoder) Encode(data []byte) ([]byte, error) {
	out := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(out, data)

	return out, nil
}

func (e hexEncoder) Decode(data []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(data)))
	_, err := hex.Decode(out, data)
	if err != nil {
		return nil, err
	}

	return out, nil
}
