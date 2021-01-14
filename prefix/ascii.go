package prefix

import (
	"fmt"
	"strconv"
	"strings"
)

type asciiVarPrefixer struct {
	Digits int
}

var ASCII = Prefixers{
	Fixed: &asciiFixedPrefixer{},
	L:     &asciiVarPrefixer{1},
	LL:    &asciiVarPrefixer{2},
	LLL:   &asciiVarPrefixer{3},
	LLLL:  &asciiVarPrefixer{4},
}

func (p *asciiVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d is larger than maximum: %d", dataLen, maxLen)
	}

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, fmt.Errorf("Failed to encode length: %d. Number of digits exceeds: %d", dataLen, p.Digits)
	}

	res := fmt.Sprintf("%0*d", p.Digits, dataLen)

	return []byte(res), nil
}

func (p *asciiVarPrefixer) DecodeLength(maxLen int, data []byte) (int, error) {
	if len(data) < p.Digits {
		return 0, fmt.Errorf("Failed to decode length. Not enought data length: %d to read: %d byte digits", len(data), p.Digits)
	}

	dataLen, err := strconv.Atoi(string(data[:p.Digits]))
	if err != nil {
		return 0, fmt.Errorf("Failed to decode length: %w", err)
	}

	if dataLen > maxLen {
		return 0, fmt.Errorf("Failed to decode length. Data length %d is larger than maximum %d", dataLen, maxLen)
	}

	return dataLen, nil
}

func (p *asciiVarPrefixer) Length() int {
	return p.Digits
}

func (p *asciiVarPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII %s length", strings.Repeat("L", p.Digits))
}

type asciiFixedPrefixer struct {
}

func (p *asciiFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d should be fixed: %d", dataLen, fixLen)
	}

	return []byte{}, nil
}

func (p *asciiFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return fixLen, nil
}

func (p *asciiFixedPrefixer) Length() int {
	return 0
}

func (p *asciiFixedPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII fixed length")
}
