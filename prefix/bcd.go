package prefix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/moov-io/iso8583/encoding"
	"github.com/yerden/go-util/bcd"
)

type bcdVarPrefixer struct {
	Digits int
}

var BCD = Prefixers{
	Fixed:  &bcdFixedPrefixer{},
	L:      &bcdVarPrefixer{1},
	LL:     &bcdVarPrefixer{2},
	LLL:    &bcdVarPrefixer{3},
	LLLL:   &bcdVarPrefixer{4},
	LLLLL:  &bcdVarPrefixer{5},
	LLLLLL: &bcdVarPrefixer{6},
}

func (p *bcdVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
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

	strLen := fmt.Sprintf("%0*d", p.Digits, dataLen)
	res, err := encoding.BCD.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *bcdVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	length := bcd.EncodedLen(p.Digits)
	if len(data) < length {
		return 0, 0, &LengthError{
			fmt.Errorf(notEnoughDataToRead, length, len(data)),
		}
	}

	bDigits, _, err := encoding.BCD.Decode(data[:length], p.Digits)
	if err != nil {
		return 0, 0, err
	}

	dataLen, err := strconv.Atoi(string(bDigits))
	if err != nil {
		return 0, 0, err
	}

	if dataLen > maxLen {
		return 0, 0, &LengthError{
			fmt.Errorf(dataLengthIsLargerThanMax, dataLen, maxLen),
		}
	}

	return dataLen, length, nil
}

func (p *bcdVarPrefixer) Inspect() string {
	return fmt.Sprintf("BCD.%s", strings.Repeat("L", p.Digits))
}

type bcdFixedPrefixer struct {
}

func (p *bcdFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, &LengthError{
			fmt.Errorf(fieldLengthShouldBeFixed, dataLen, fixLen),
		}
	}

	return []byte{}, nil
}

// Returns number of characters that should be decoded
func (p *bcdFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *bcdFixedPrefixer) Inspect() string {
	return "BCD.Fixed"
}
