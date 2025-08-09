package encoding

import (
	"bytes"
	"fmt"

	"github.com/yerden/go-util/bcd"

	"github.com/moov-io/iso8583/utils"
)

// trackDigitsBCD defines the BCD mapping for digits '0'-'9' and the separator 'D'.
var trackDigitsBCD = &bcd.BCD{
	Map: func() map[byte]byte {
		m := make(map[byte]byte, 11)
		for c := byte('0'); c <= byte('9'); c++ {
			m[c] = c - '0'
		}
		m['D'] = 0xD
		return m
	}(),
	Filler:      0xF,
	SwapNibbles: false,
}

var (
	_         Encoder = (*track2BcdEncoder)(nil)
	Track2BCD         = &track2BcdEncoder{}
)

type track2BcdEncoder struct{}

// Encode converts a byte slice into BCD format.
func (e track2BcdEncoder) Encode(src []byte) ([]byte, error) {
	// if the number of digits is odd, add ‘0’ to the left
	if len(src)%2 != 0 {
		src = append([]byte("0"), src...)
	}

	enc := bcd.NewEncoder(trackDigitsBCD)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to perform BCD encoding")
	}
	return dst[:n], nil
}

// Decode converts a BCD-encoded byte slice back to its original representation.
func (e track2BcdEncoder) Decode(src []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	// length is in nibbles; bytes to read = ceil(n/2)
	decodedLen := length
	if decodedLen%2 != 0 {
		decodedLen++
	}
	read := bcd.EncodedLen(decodedLen) // = (decodedLen+1) / 2

	if len(src) < read {
		return nil, 0, fmt.Errorf("not enough data to decode. expected len %d, got %d", read, len(src))
	}

	dec := bcd.NewDecoder(trackDigitsBCD)
	dst := make([]byte, decodedLen)
	_, err := dec.Decode(dst, src[:read])
	if err != nil {
		return nil, 0, utils.NewSafeError(err, "failed to perform BCD decoding")
	}

	// remove the padding '0' if n is odd and truncate to 'length' nibbles.
	out := dst[decodedLen-length:]

	// reject hex nibbles outside the allowed set ('D').
	invalidNibbleIndex := bytes.IndexFunc(out, func(r rune) bool {
		return r >= 'A' && r <= 'F' && r != 'D'
	})
	if invalidNibbleIndex != -1 {
		return nil, 0, fmt.Errorf("invalid track2 nibble decoded: %q", out[invalidNibbleIndex])
	}
	return out, read, nil
}
