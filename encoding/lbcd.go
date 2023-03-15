package encoding

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
	"github.com/yerden/go-util/bcd"
)

var (
	_    Encoder = (*lBCDEncoder)(nil)
	LBCD         = &lBCDEncoder{}
)

type lBCDEncoder struct{}

func (e *lBCDEncoder) Encode(src []byte) ([]byte, error) {
	if len(src)%2 != 0 {
		src = append(src, []byte("0")...)
	}

	enc := bcd.NewEncoder(bcd.Standard)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to perform BCD encoding")
	}

	return dst[:n], nil
}

func (e *lBCDEncoder) Decode(src []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	decodedLen := length
	if length%2 != 0 {
		decodedLen += 1
	}

	read := bcd.EncodedLen(decodedLen)

	dec := bcd.NewDecoder(bcd.Standard)
	dst := make([]byte, decodedLen)

	if len(src) < read {
		return nil, 0, fmt.Errorf("not enough data to decode. expected len %d, got %d", read, len(src))
	}

	_, err := dec.Decode(dst, src)
	if err != nil {
		return nil, 0, utils.NewSafeError(err, "failed to perform BCD decoding")
	}

	// because it's left aligned, we return data from
	// 0 index
	return dst[:length], read, nil
}
