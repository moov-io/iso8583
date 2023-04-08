package encoding

import (
	"fmt"

	"github.com/moov-io/iso8583/utils"
	"github.com/yerden/go-util/bcd"
)

var (
	_   Encoder = (*bcdEncoder)(nil)
	BCD         = &bcdEncoder{}
)

type bcdEncoder struct{}

// Encode returns packed data and the original length in digits
func (e *bcdEncoder) Encode(src []byte) ([]byte, int, error) {
	length := len(src)
	if len(src)%2 != 0 {
		src = append([]byte("0"), src...)
	}

	enc := bcd.NewEncoder(bcd.Standard)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, 0, utils.NewSafeError(err, "failed to perform BCD encoding")
	}

	return dst[:n], length, nil
}

// Decode returns bcd unpacked data and number of digits decoded
func (e *bcdEncoder) Decode(src []byte, length int) ([]byte, int, error) {
	// length should be positive
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	// for BCD encoding the length should be even
	decodedLen := length
	if length%2 != 0 {
		decodedLen += 1
	}

	// how many bytes we will read
	read := bcd.EncodedLen(decodedLen)

	if len(src) < read {
		return nil, 0, fmt.Errorf("not enough data to decode. expected len %d, got %d", read, len(src))
	}

	dec := bcd.NewDecoder(bcd.Standard)
	dst := make([]byte, decodedLen)
	_, err := dec.Decode(dst, src[:read])
	if err != nil {
		return nil, 0, utils.NewSafeError(err, "failed to perform BCD decoding")
	}

	// becase BCD is right aligned, we skip first bytes and
	// read only what we need
	// e.g. 0643 => 643
	return dst[decodedLen-length:], read, nil
}
