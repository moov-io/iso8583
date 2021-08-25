package encoding

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// ASCII HEX encoder
var Hex = &hexEncoder{}

type hexEncoder struct{}

func (e hexEncoder) Encode(data []byte) ([]byte, error) {
	out := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(out, data)

	str := string(out)
	str = strings.ToUpper(str)

	return []byte(str), nil
}

// Decodes ASCII hex and returns bytes
// length is number of HEX-digits (two ASCII characters is one HEX digit)
func (e hexEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	// to read 8 HEX digits we have to read 16 ASCII chars (bytes)
	read := hex.EncodedLen(length)
	if read > len(data) {
		return nil, 0, fmt.Errorf("not enough data to read")
	}

	out := make([]byte, length)

	_, err := hex.Decode(out, data[:read])
	if err != nil {
		return nil, 0, err
	}

	return out, read, nil
}
