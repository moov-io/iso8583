package encoding

import (
	"github.com/yerden/go-util/bcd"
)

var LBCD Encoder = &lBCDEncoder{}

type lBCDEncoder struct{}

func (e *lBCDEncoder) Encode(src []byte) ([]byte, error) {
	if len(src)%2 != 0 {
		src = append(src, []byte("0")...)
	}

	enc := bcd.NewEncoder(bcd.Standard)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

func (e *lBCDEncoder) Decode(src []byte, length int) ([]byte, int, error) {
	decodedLen := length
	if length%2 != 0 {
		decodedLen += 1
	}

	read := bcd.EncodedLen(decodedLen)

	dec := bcd.NewDecoder(bcd.Standard)
	dst := make([]byte, decodedLen)
	_, err := dec.Decode(dst, src)
	if err != nil {
		return nil, 0, err
	}

	// because it's left aligned, we return data from
	// 0 index
	return dst[:length], read, nil
}
