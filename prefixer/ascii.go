package prefixer

import (
	"fmt"
	"strconv"
	"strings"
)

type asciiVarPrefixer struct {
	MaxLen int
	Digits int
}

var ASCII = Prefixers{
	Fixed: func(fixLen int) Prefixer { return &asciiFixedPrefixer{fixLen} },
	L:     func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 1} },
	LL:    func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 2} },
	LLL:   func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 3} },
	LLLL:  func(maxLen int) Prefixer { return &asciiVarPrefixer{maxLen, 4} },
}

func (p *asciiVarPrefixer) EncodeLength(dataLen int) ([]byte, error) {
	if dataLen > p.MaxLen {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d is larger than maximum: %d", dataLen, p.MaxLen)
	}

	// convert int into []byte
	res := strconv.AppendInt([]byte{}, int64(dataLen), 10)
	if len(res) > p.Digits {
		return nil, fmt.Errorf("Failed to encode length: %d. Number of digits exceeds: %d", dataLen, p.Digits)
	}

	return res, nil
}

func (p *asciiVarPrefixer) DecodeLength(data []byte) (int, error) {
	if len(data) < p.Digits {
		return 0, fmt.Errorf("Failed to decode length. Not enought data length: %d to read: %d byte digits", len(data), p.Digits)
	}

	dataLen, err := strconv.Atoi(string(data[:p.Digits]))
	if err != nil {
		return 0, fmt.Errorf("Failed to decode length: %w", err)
	}

	if dataLen > p.MaxLen {
		return 0, fmt.Errorf("Failed to decode length. Data length %d is larger than maximum %d", dataLen, p.MaxLen)
	}

	return dataLen, nil
}

func (p *asciiVarPrefixer) Length() int {
	return p.Digits
}

func (p *asciiVarPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII %s prefixer. Max Length: %d", strings.Repeat("L", p.Digits), p.MaxLen)
}

type asciiFixedPrefixer struct {
	Len int
}

func (p *asciiFixedPrefixer) EncodeLength(dataLen int) ([]byte, error) {
	if dataLen != p.Len {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d should be fixed: %d", dataLen, p.Len)
	}

	return []byte{}, nil
}

func (p *asciiFixedPrefixer) DecodeLength(data []byte) (int, error) {
	return p.Len, nil
}

func (p *asciiFixedPrefixer) Length() int {
	return 0
}

func (p *asciiFixedPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII fixed prefixer. Length: %d", p.Len)
}
