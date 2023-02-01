package prefix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

var Binary = Prefixers{
	Fixed: &binaryFixedPrefixer{},
	L:     &binaryVarPrefixer{1},
	LL:    &binaryVarPrefixer{2},
	LLL:   &binaryVarPrefixer{3},
	LLLL:  &binaryVarPrefixer{4},
}

type binaryFixedPrefixer struct {
}

func (p *binaryFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen != fixLen {
		return nil, fmt.Errorf("field length: %d should be fixed: %d", dataLen, fixLen)
	}

	return []byte{}, nil
}

func (p *binaryFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *binaryFixedPrefixer) Inspect() string {
	return "Binary.Fixed"
}

type binaryVarPrefixer struct {
	Digits int
}

func intToBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("negative number: %d", n)
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint32(n))
	if err != nil {
		return nil, fmt.Errorf("int to bytes: %w", err)
	}
	return buf.Bytes(), nil
}

func bytesToInt(b []byte) (int, error) {
	buf := bytes.NewReader(b)
	var n uint32
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		return 0, fmt.Errorf("bytes to int: %w", err)
	}
	return int(n), nil
}

func (p *binaryVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("field length: %d is larger than maximum: %d", dataLen, maxLen)
	}

	res, err := intToBytes(dataLen)
	if err != nil {
		return nil, fmt.Errorf("encode length: %w", err)
	}

	// remove all leading zeros as res is always 4 bytes
	res = bytes.TrimLeft(res, "\x00")

	if len(res) > p.Digits {
		return nil, fmt.Errorf("number of digits in length: %d exceeds: %d", dataLen, p.Digits)
	}

	// if len of res is less than p.Digits prepend with 0x00
	if len(res) < p.Digits {
		res = append(bytes.Repeat([]byte{0x00}, p.Digits-len(res)), res...)
	}

	return res, nil
}

// DecodeLength decodes the length of the field from the data. It reads up to 4
// bytes from data, converts it into int32 and returns the length of the field
// and the number of bytes read.
func (p *binaryVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	if len(data) < p.Digits {
		return 0, 0, fmt.Errorf("not enough data length: %d to read: %d bytes", len(data), p.Digits)
	}

	prefBytes := data[:p.Digits]

	// it take 4 bytes to encode (u)int32
	uint32Size := 4

	// prepend with 0x00 if len of data is less than intSize (4 bytes)
	if len(prefBytes) < uint32Size {
		prefBytes = append(bytes.Repeat([]byte{0x00}, uint32Size-len(prefBytes)), prefBytes...)
	}

	dataLen, err := bytesToInt(prefBytes)
	if err != nil {
		return 0, 0, fmt.Errorf("decode length: %w", err)
	}

	if dataLen > maxLen {
		return 0, 0, fmt.Errorf("data length: %d is larger than maximum %d", dataLen, maxLen)
	}

	return dataLen, p.Digits, nil
}

func (p *binaryVarPrefixer) Inspect() string {
	return fmt.Sprintf("Binary.%s", strings.Repeat("L", p.Digits))
}
