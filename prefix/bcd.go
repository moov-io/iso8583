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
	Fixed: &bcdFixedPrefixer{},
	L:     &bcdVarPrefixer{1},
	LL:    &bcdVarPrefixer{2},
	LLL:   &bcdVarPrefixer{3},
	LLLL:  &bcdVarPrefixer{4},
}

func (p *bcdVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("Failed to encode length. Field length: %d is larger than maximum: %d", dataLen, maxLen)
	}

	if len(strconv.Itoa(dataLen)) > p.Digits {
		return nil, fmt.Errorf("Failed to encode length: %d. Number of digits exceeds: %d", dataLen, p.Digits)
	}

	strLen := fmt.Sprintf("%0*d", p.Digits, dataLen)
	res, err := encoding.BCD.Encode([]byte(strLen))

	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	return res, nil
}

func (p *bcdVarPrefixer) DecodeLength(maxLen int, data []byte) (int, error) {
	if len(data) < p.Length() {
		return 0, fmt.Errorf("length mismatch: want to read %d bytes, get only %d", p.Length(), len(data))
	}

	bDigits, err := encoding.BCD.Decode(data[:p.Length()])
	if err != nil {
		return 0, fmt.Errorf("failed to decode BCD encoded length: %w", err)
	}

	dataLen, err := strconv.Atoi(string(bDigits[len(bDigits)-p.Digits:]))
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}

	if dataLen > maxLen {
		return 0, fmt.Errorf("failed to decode length. Data length %d is larger than maximum %d", dataLen, maxLen)
	}

	return dataLen, nil
}

func (p *bcdVarPrefixer) Length() int {
	return bcd.EncodedLen(p.Digits)
}

func (p *bcdVarPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII %s length", strings.Repeat("L", p.Digits))
}

type bcdFixedPrefixer struct {
}

func (p *bcdFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen > fixLen {
		return nil, fmt.Errorf("failed to encode length: field length: %d should be fixed: %d", dataLen, fixLen)
	}

	return []byte{}, nil
}

func (p *bcdFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, error) {
	return bcd.EncodedLen(fixLen), nil
}

func (p *bcdFixedPrefixer) Length() int {
	return 0
}

func (p *bcdFixedPrefixer) Inspect() string {
	return fmt.Sprintf("ASCII fixed length")
}
