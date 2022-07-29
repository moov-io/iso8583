package prefix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/moov-io/iso8583/encoding"
)

var EBCDIC1047 = Prefixers{
	Fixed: &ebcdic1047FixedPrefixer{},
	L:     &ebcdic1047Prefixer{1},
	LL:    &ebcdic1047Prefixer{2},
	LLL:   &ebcdic1047Prefixer{3},
	LLLL:  &ebcdic1047Prefixer{4},
}

type ebcdic1047Prefixer struct {
	digits int
}

func (p *ebcdic1047Prefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("field length [%d] is larger than maximum [%d]", dataLen, maxLen)
	}

	if len(strconv.Itoa(dataLen)) > p.digits {
		return nil, fmt.Errorf("number of digits in data [%d] exceeds its maximum indicator [%d]", dataLen, p.digits)
	}

	strLen := fmt.Sprintf("%0*d", p.digits, dataLen)
	res, err := encoding.EBCDIC1047.Encode([]byte(strLen))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ebcdic1047Prefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	if len(data) < p.digits {
		return 0, 0, fmt.Errorf("not enough data length [%d] to read [%d] byte digits", len(data), p.digits)
	}

	decodedData, _, err := encoding.EBCDIC1047.Decode(data[:p.digits], p.digits)
	if err != nil {
		return 0, 0, err
	}

	dataLen, err := strconv.Atoi(string(decodedData))
	if err != nil {
		return 0, 0, fmt.Errorf("length [%s] is not a valid integer length field", string(decodedData))
	}

	if dataLen > maxLen {
		return 0, 0, fmt.Errorf("data length [%d] is larger than maximum [%d]", dataLen, maxLen)
	}

	return dataLen, p.digits, nil
}

func (p *ebcdic1047Prefixer) Inspect() string {
	return fmt.Sprintf("EBCDIC.%s", strings.Repeat("L", p.digits))
}

type ebcdic1047FixedPrefixer struct{}

func (p *ebcdic1047FixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, fmt.Errorf("field length [%d] should be fixed [%d]", dataLen, fixLen)
	}

	return []byte{}, nil
}

func (p *ebcdic1047FixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *ebcdic1047FixedPrefixer) Inspect() string {
	return "EBCDIC.Fixed"
}
