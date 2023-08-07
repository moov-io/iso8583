package prefix

import (
	"fmt"
	"strconv"
	"strings"
)

type TLVPrefixer struct {
	Digits int
}

var TLV = &Prefixers{
	L:    &TLVPrefixer{1},
	LL:   &TLVPrefixer{2},
	LLL:  &TLVPrefixer{3},
	LLLL: &TLVPrefixer{4},
}

func (p *TLVPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("field length: %d is larger than maximum: %d", dataLen, maxLen)
	}

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, fmt.Errorf("number of digits in length: %d exceeds: %d", dataLen, p.Digits)
	}

	res := fmt.Sprintf("%0*d", p.Digits, dataLen)

	return []byte(res), nil
}

func (p *TLVPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	if len(data) < p.Digits {
		return 0, 0, fmt.Errorf("not enough data length: %d to read: %d byte digits", len(data), p.Digits)
	}

	dataLen, err := strconv.Atoi(string(data[:p.Digits]))
	if err != nil {
		return 0, 0, err
	}

	// length should be positive
	if dataLen < 0 {
		return 0, 0, fmt.Errorf("invalid length: %d", dataLen)
	}

	if dataLen > maxLen {
		return 0, 0, fmt.Errorf("data length: %d is larger than maximum %d", dataLen, maxLen)
	}

	return dataLen, p.Digits, nil
}

func (p *TLVPrefixer) Inspect() string {
	return fmt.Sprintf("TLV.%s", strings.Repeat("L", p.Digits))
}
