package prefix

import (
	"fmt"
)

var Binary = Prefixers{
	Fixed: &binaryFixedPrefixer{},
}

type binaryFixedPrefixer struct {
}

func (p *binaryFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, fmt.Errorf("field length: %d should be fixed: %d", dataLen, fixLen)
	}

	return []byte{}, nil
}

func (p *binaryFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *binaryFixedPrefixer) Inspect() string {
	return "Binary.Fixed"
}
