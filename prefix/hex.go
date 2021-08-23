package prefix

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/franizus/iso8583/encoding"
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

func (p *hexFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return fixLen, nil
}

func (p *hexFixedPrefixer) Length() int {
	return 0
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

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, fmt.Errorf("number of digits in length: %d exceeds: %d", dataLen, p.Digits)
	}

	strLen := strconv.Itoa(dataLen)
	res, err := encoding.Hex.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *hexVarPrefixer) DecodeLength(maxLen int, data []byte) (int, error) {
	if len(data) < p.Length() {
		return 0, fmt.Errorf("length mismatch: want to read %d bytes, get only %d", p.Length(), len(data))
	}

	bDigits, _, err := encoding.Hex.Decode(data[:p.Length()], p.Digits)
	if err != nil {
		return 0, err
	}

	dataLen, err := strconv.Atoi(string(bDigits))
	if err != nil {
		return 0, err
	}

	if dataLen > maxLen {
		return 0, fmt.Errorf("data length %d is larger than maximum %d", dataLen, maxLen)
	}

	return dataLen, nil
}

func (p *hexVarPrefixer) Length() int {
	return hex.EncodedLen(p.Digits)
}

func (p *hexVarPrefixer) Inspect() string {
	return fmt.Sprintf("Hex.%s", strings.Repeat("L", p.Digits))
}
