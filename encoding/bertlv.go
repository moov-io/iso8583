package encoding

import (
        "bytes"
	"fmt"
)

// BER-TLV Tag encoder
var BerTLVTag Encoder = &berTLVEncoderTag{}

type berTLVEncoderTag struct{}

// Encode converts ASCII Hex-digits into a byte slice
// e.g. []byte("AABBCC") would be converted into
// []byte{0xAA, 0xBB, 0xCC}
func (berTLVEncoderTag) Encode(data []byte) ([]byte, error) {
        out, _, err := Hex.Decode(data, len(data))
        return out, err
}

// Decode converts hexadecimal TLV bytes into their ASCII representation. 
func (berTLVEncoderTag) Decode(data []byte, length int) ([]byte, int, error) {
        r := bytes.NewReader(data)

	firstByte, err := r.ReadByte()
	if err != nil {
                return nil, 0, err
	}
        tagLen := 1

        // If bits 5 - 1 of the tag's first byte are all set, then we must read
        // the subsequent byte for the tag number. We can inspect this by
        // simply checking whether the first byte of the tag >= 0b00011111
        shouldReadSubsequentByte := true
	if firstByte >= 0b00011111 {
                shouldReadSubsequentByte = false
	}

	for shouldReadSubsequentByte {
		b, err := r.ReadByte()
		if err != nil {
                        return nil, tagLen, fmt.Errorf("failed to decode TLV tag: %w", err)
		}
		tagLen++
                // We read subsequent bytes to extract the tag by checking if
                // the the most significant bit is set.
		if (b >> 7) == 0 {
			shouldReadSubsequentByte = false
		}
	}

        out, err := Hex.Encode(data[:tagLen])
	if err != nil {
                return nil, tagLen, err
	}
        return out, tagLen, nil
}
