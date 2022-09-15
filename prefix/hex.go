package prefix

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

var Hex = Prefixers{
	Fixed: &hexFixedPrefixer{},
	L:     &hexVarPrefixer{1},
	LL:    &hexVarPrefixer{2},
	LLL:   &hexVarPrefixer{3},
	LLLL:  &hexVarPrefixer{4},
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

func (p *hexFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *hexFixedPrefixer) Inspect() string {
	return "Hex.Fixed"
}

type hexVarPrefixer struct {
	Digits int
}

func (p *hexVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("field length: %d is larger than maximum: %d", dataLen, maxLen)
	}

	maxPossibleLength := 1<<(p.Digits*8) - 1
	if dataLen > maxPossibleLength {
		return nil, fmt.Errorf("number of digits in length: %d exceeds: %d", dataLen, p.Digits)
	}

	strLen := strconv.FormatInt(int64(dataLen), 16)
	res := fmt.Sprintf("%0*s", p.Digits*2, strings.ToUpper(strLen))

	return []byte(res), nil
}

func (p *hexVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	length := hex.EncodedLen(p.Digits)
	if len(data) < length {
		return 0, 0, fmt.Errorf("length mismatch: want to read %d bytes, get only %d", length, len(data))
	}

	dataLen, err := strconv.ParseInt(string(data[:length]), 16, p.Digits*8)
	if err != nil {
		return 0, 0, err
	}

	if int(dataLen) > maxLen {
		return 0, 0, fmt.Errorf("data length %d is larger than maximum %d", dataLen, maxLen)
	}

	return int(dataLen), length, nil
}

func (p *hexVarPrefixer) Inspect() string {
	return fmt.Sprintf("Hex.%s", strings.Repeat("L", p.Digits))
}
