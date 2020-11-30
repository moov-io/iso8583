package prefixer

import "fmt"

var Hex = Prefixers{
	Fixed: func(fixLen int) Prefixer { return &hexFixedPrefixer{fixLen} },

	// return not implemented error
	L:    func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 1} },
	LL:   func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 2} },
	LLL:  func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 3} },
	LLLL: func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 4} },
}

type hexFixedPrefixer struct {
	Len int
}

func (p *hexFixedPrefixer) EncodeLength(dataLen int) ([]byte, error) {
	// for ascii hex the lenght is x2 (ascii hex digit takes one byte)
	if dataLen != p.Len*2 {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d should be fixed: %d", dataLen, p.Len*2)
	}

	return []byte{}, nil
}

func (p *hexFixedPrefixer) DecodeLength(data []byte) (int, error) {
	return p.Len * 2, nil
}

func (p *hexFixedPrefixer) Length() int {
	return 0
}

func (p *hexFixedPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII fixed prefixer. Length: %d", p.Len)
}
