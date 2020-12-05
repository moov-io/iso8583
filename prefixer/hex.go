package prefixer

import "fmt"

var Hex = Prefixers{
	Fixed: &hexFixedPrefixer{},
}

type hexFixedPrefixer struct {
}

func (p *hexFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	// for ascii hex the lenght is x2 (ascii hex digit takes one byte)
	if dataLen != fixLen*2 {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d should be fixed: %d", dataLen, fixLen*2)
	}

	return []byte{}, nil
}

func (p *hexFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return fixLen * 2, nil
}

func (p *hexFixedPrefixer) Length() int {
	return 0
}

func (p *hexFixedPrefixer) Inspect() string {
	return "ASCII fixed length"
}
