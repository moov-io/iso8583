package prefix

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var Hex = Prefixers{
	Fixed:  &hexFixedPrefixer{},
	L:      &hexVarPrefixer{1},
	LL:     &hexVarPrefixer{2},
	LLL:    &hexVarPrefixer{3},
	LLLL:   &hexVarPrefixer{4},
	LLLLL:  &hexVarPrefixer{5},
	LLLLLL: &hexVarPrefixer{6},
}

type hexFixedPrefixer struct {
}

func (p *hexFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	// for ascii hex the length is x2 (ascii hex digit takes one byte)
	if dataLen != fixLen*2 {
		return nil, fmt.Errorf(fieldLengthShouldBeFixed, dataLen, fixLen*2)
	}

	return []byte{}, nil
}

func (p *hexFixedPrefixer) DecodeLength(fixLen int, _ []byte) (int, int, error) {
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
		return nil, fmt.Errorf(fieldLengthIsLargerThanMax, dataLen, maxLen)
	}

	maxPossibleLength := 1<<(p.Digits*8) - 1
	if dataLen > maxPossibleLength {
		return nil, fmt.Errorf(numberOfDigitsInLengthExceeds, dataLen, p.Digits)
	}

	strLen := strconv.FormatInt(int64(dataLen), 16)
	res := fmt.Sprintf("%0*s", p.Digits*2, strings.ToUpper(strLen))

	return []byte(res), nil
}

func (p *hexVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	length := hex.EncodedLen(p.Digits)
	if len(data) < length {
		return 0, 0, fmt.Errorf(notEnoughDataToRead, length, len(data))
	}

	dataLen, err := strconv.ParseUint(string(data[:length]), 16, p.Digits*8)
	if err != nil {
		return 0, 0, err
	}

	if dataLen > math.MaxInt {
		return 0, 0, fmt.Errorf("data length %d exceeds maximum int value", dataLen)
	}

	// #nosec G115 -- dataLen is validated to be within MaxInt range above
	if int(dataLen) > maxLen {
		return 0, 0, fmt.Errorf(dataLengthIsLargerThanMax, dataLen, maxLen)
	}

	// #nosec G115 -- dataLen is validated to be within MaxInt range above
	return int(dataLen), length, nil
}

func (p *hexVarPrefixer) Inspect() string {
	return fmt.Sprintf("Hex.%s", strings.Repeat("L", p.Digits))
}
