package encoding

import (
	"github.com/yerden/go-util/bcd"
)

var BCD Encoder = &bcdEncoder{}

type bcdEncoder struct{}

func (e *bcdEncoder) Encode(src []byte) ([]byte, error) {
	enc := bcd.NewEncoder(bcd.Standard)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}

func (e *bcdEncoder) Decode(src []byte) ([]byte, error) {
	dec := bcd.NewDecoder(bcd.Standard)
	dst := make([]byte, bcd.DecodedLen(len(src)))
	n, err := dec.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	return dst[:n], nil
}
