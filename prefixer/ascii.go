package prefixer

import (
	"fmt"
	"strconv"
	"strings"
)

type asciiPrefixer struct {
	MaxLen int
	Digits int
}

var ASCII = varPrefixers{
	L:    func(maxLen int) Prefixer { return &asciiPrefixer{maxLen, 1} },
	LL:   func(maxLen int) Prefixer { return &asciiPrefixer{maxLen, 2} },
	LLL:  func(maxLen int) Prefixer { return &asciiPrefixer{maxLen, 3} },
	LLLL: func(maxLen int) Prefixer { return &asciiPrefixer{maxLen, 4} },
}

func (p *asciiPrefixer) EncodeLength(dataLen int) ([]byte, error) {
	if dataLen > p.MaxLen {
		return nil, fmt.Errorf("Failed to encode length. Provided length: %d is larger than maximum: %d", dataLen, p.MaxLen)
	}

	// convert int into []byte
	res := strconv.AppendInt([]byte{}, int64(dataLen), 10)
	if len(res) > p.Digits {
		return nil, fmt.Errorf("Failed to encode length: %d. Number of digits exceeds: %d", dataLen, p.Digits)
	}

	return res, nil
}

func (p *asciiPrefixer) DecodeLength(data []byte) (int, error) {
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

func (p *asciiPrefixer) DecodedLength() int {
	return p.Digits
}

func (p *asciiPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII %s prefixer. Max Length: %d", strings.Repeat("L", p.Digits), p.MaxLen)
}
