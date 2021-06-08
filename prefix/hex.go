package prefix

import "fmt"

var Hex = Prefixers{
	Fixed: &hexFixedPrefixer{},
}

type hexFixedPrefixer struct {
}

func (p *hexFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	// for ascii hex the length is x2 (ascii hex digit takes one byte)
	if dataLen != fixLen*2 {
		return nil, fmt.Errorf("field length: %d should be fixed: %d", dataLen, fixLen*2)
	}

	return []byte{}, nil
}

func (p *hexFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return fixLen, nil
}

func (p *hexFixedPrefixer) Length() int {
	return 0
}

func (p *hexFixedPrefixer) Inspect() string {
	return "Hex fixed length"
}
