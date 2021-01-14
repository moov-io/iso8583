package encoding

import "fmt"

var ASCII Encoder = &asciiEncoder{}

type asciiEncoder struct{}

func (e asciiEncoder) Encode(data []byte) ([]byte, error) {
	out := []byte{}
	for _, r := range data {
		if r > 127 {
			return nil, fmt.Errorf("Failed to decode into ASCII. '%s' is not valid ASCII char", string(r))
		}
		out = append(out, r)
	}

	return out, nil
}

func (e asciiEncoder) Decode(data []byte, _ int) ([]byte, error) {
	out := []byte{}
	for _, r := range data {
		if r > 127 {
			return nil, fmt.Errorf("Failed to decode into ASCII. '%s' is not valid ASCII char", string(r))
		}
		out = append(out, r)
	}

	return out, nil
}
