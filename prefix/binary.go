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

func (p *binaryFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return fixLen, nil
}

func (p *binaryFixedPrefixer) Length() int {
	return 0
}

func (p *binaryFixedPrefixer) Inspect() string {
	return "Binary fixed length"
}

func (p *binaryFixedPrefixer) InspectName() string {
	return "Fixed"
}
