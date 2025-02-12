package prefix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/moov-io/iso8583/encoding"
)

type ebcdicVarPrefixer struct {
	Digits int
}

var EBCDIC = Prefixers{
	Fixed:  &ebcdicFixedPrefixer{},
	L:      &ebcdicVarPrefixer{1},
	LL:     &ebcdicVarPrefixer{2},
	LLL:    &ebcdicVarPrefixer{3},
	LLLL:   &ebcdicVarPrefixer{4},
	LLLLL:  &ebcdicVarPrefixer{5},
	LLLLLL: &ebcdicVarPrefixer{6},
}

func (p *ebcdicVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf(fieldLengthIsLargerThanMax, dataLen, maxLen)
	}

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, fmt.Errorf(numberOfDigitsInLengthExceeds, dataLen, p.Digits)
	}

	strLen := fmt.Sprintf("%0*d", p.Digits, dataLen)
	res, err := encoding.EBCDIC.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *ebcdicVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	length := p.Digits
	if len(data) < length {
		return 0, 0, fmt.Errorf(notEnoughDataToRead, length, len(data))
	}

	bDigits, _, err := encoding.EBCDIC.Decode(data[:length], p.Digits)
	if err != nil {
		return 0, 0, err
	}

	dataLen, err := strconv.Atoi(string(bDigits))
	if err != nil {
		return 0, 0, err
	}

	if dataLen > maxLen {
		return 0, 0, fmt.Errorf(dataLengthIsLargerThanMax, dataLen, maxLen)
	}

	return dataLen, length, nil
}

func (p *ebcdicVarPrefixer) Inspect() string {
	return fmt.Sprintf("EBCDIC.%s", strings.Repeat("L", p.Digits))
}

type ebcdicFixedPrefixer struct {
}

func (p *ebcdicFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, fmt.Errorf(fieldLengthShouldBeFixed, dataLen, fixLen)
	}

	return []byte{}, nil
}

// Returns number of characters that should be decoded
func (p *ebcdicFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *ebcdicFixedPrefixer) Inspect() string {
	return "EBCDIC.Fixed"
}
