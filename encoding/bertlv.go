package encoding

import (
        "bytes"
        "encoding/hex"
	"fmt"
	"math/bits"
)

// BER-TLV Tag encoder
var BerTLVTag Encoder = &berTLVEncoderTag{}

type berTLVEncoderTag struct{}

// Encode converts ASCII Hex-digits into a byte slice e.g. []byte("AABBCC")
// would be converted into []byte{0xAA, 0xBB, 0xCC}
func (berTLVEncoderTag) Encode(data []byte) ([]byte, error) {
        out, _, err := Hex.Decode(data, hex.DecodedLen(len(data)))
        return out, err
}

// Decode converts hexadecimal TLV bytes into their ASCII representation according
// to the following rules:
//
// 1) If bits 5 - 1 of the tag's first byte are all set, then we must read
//    the subsequent byte for the tag number.
// 2) We must continue reading subsequent bytes until we arrive at one whose
//    most significant bit is unset.
//
// On success, the ASCII representation of the Tag as well are returned along
// with the number of bytes read e.g. []byte{0x5F, 0x2A} would be converted to
// []byte("5F2A")
func (berTLVEncoderTag) Decode(data []byte, length int) ([]byte, int, error) {
        r := bytes.NewReader(data)

	firstByte, err := r.ReadByte()
	if err != nil {
                return nil, 0, err
	}
        tagLen := 1

        shouldReadSubsequentByte := false
        if bits.TrailingZeros8(^firstByte) >= 5 {
                shouldReadSubsequentByte = true
	}

	for shouldReadSubsequentByte {
		b, err := r.ReadByte()
		if err != nil {
                        return nil, tagLen, fmt.Errorf("failed to decode TLV tag: %w", err)
		}
		tagLen++
                // We read subsequent bytes to extract the tag by checking if
                // the the most significant bit is set.
                if bits.LeadingZeros8(b) > 0 {
			shouldReadSubsequentByte = false
		}
	}

        out, err := Hex.Encode(data[:tagLen])
	if err != nil {
                return nil, tagLen, err
	}
        return out, tagLen, nil
}
