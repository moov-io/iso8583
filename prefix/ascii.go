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
	Fixed:  &asciiFixedPrefixer{},
	L:      &asciiVarPrefixer{1},
	LL:     &asciiVarPrefixer{2},
	LLL:    &asciiVarPrefixer{3},
	LLLL:   &asciiVarPrefixer{4},
	LLLLL:  &asciiVarPrefixer{5},
	LLLLLL: &asciiVarPrefixer{6},
}

func (p *asciiVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, &LengthError{
			fmt.Errorf(fieldLengthIsLargerThanMax, dataLen, maxLen),
		}
	}

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, &LengthError{
			fmt.Errorf(numberOfDigitsInLengthExceeds, dataLen, p.Digits),
		}
	}

	res := fmt.Sprintf("%0*d", p.Digits, dataLen)

	return []byte(res), nil
}

func (p *asciiVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	if len(data) < p.Digits {
		return 0, 0, &LengthError{
			fmt.Errorf(notEnoughDataToRead, len(data), p.Digits),
		}
	}

	dataLen, err := strconv.Atoi(string(data[:p.Digits]))
	if err != nil {
		return 0, 0, err
	}

	// length should be positive
	if dataLen < 0 {
		return 0, 0, &LengthError{
			fmt.Errorf(invalidLength, dataLen),
		}
	}

	if dataLen > maxLen {
		return 0, 0, &LengthError{
			fmt.Errorf(dataLengthIsLargerThanMax, dataLen, maxLen),
		}
	}

	return dataLen, p.Digits, nil
}

func (p *asciiVarPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII.%s", strings.Repeat("L", p.Digits))
}

type asciiFixedPrefixer struct {
}

func (p *asciiFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, &LengthError{
			fmt.Errorf(fieldLengthShouldBeFixed, dataLen, fixLen),
		}
	}

	return []byte{}, nil
}

func (p *asciiFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *asciiFixedPrefixer) Inspect() string {
	return "ASCII.Fixed"
}
