package prefix

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"math/bits"
)

// BerTLV encodes and decodes the length of BER-TLV fields based on the
// following rules:
//
// Short Form: When the most-significant bit is off, the length field consists
// of only one byte in which the right-most 7 bits contain the number of bytes
// in the Value field, as an unsigned binary integer. This form of the Length
// field supports data lengths of 127 bytes.  For example, a Length value of
// 126 can be encoded as binary 01111110 (hexadecimal equivalent of 7E).
//
// Long Form: When the most-significant bit is on, the Length field consists of
// an initial byte and one or more subsequent bytes. The right-most 7 bits of
// the initial byte contain the number of subsequent bytes in the Length field,
// as an unsigned binary integer. All bits of the subsequent bytes contain an
// unsigned big-endian binary integer equal to the number of bytes in the Value
// field.  For example, a Length value of 254 can be encoded as binary 10000001
// 11111110 (hexadecimal equivalent of 81FE).
var BerTLV = &berTLVPrefixer{}

type berTLVPrefixer struct{}

// EncodeLength encodes the data length provided into a slice of bytes
// according to the rules defined above.
// NOTE: Because BER-TLV lengths are encoded dynamically, the maxLen method
// argument is ignored.
func (p *berTLVPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	buf := big.NewInt(int64(dataLen)).Bytes()
	if dataLen <= 127 {
		return buf, nil
	}
	return append([]byte{setMSB(uint8(len(buf)))}, buf...), nil
}

// DecodeLength takes in a byte array and dynamically decodes its length based
// on the rules described above. On success, both the length of the TLV value
// as well as the number bytes read to decode the length are returned.
// NOTE: Because BER-TLV lengths are decoded dynamically, the maxLen method
// argument is ignored.
func (p *berTLVPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	r := bytes.NewReader(data)

	firstByte, err := r.ReadByte()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode TLV length: %w", err)
	}

	read := 1
	if bits.LeadingZeros8(firstByte) > 0 {
		return int(firstByte), read, nil
	}

	length := make([]byte, clearMSB(firstByte))
	_, err = io.ReadFull(r, length)
	if err != nil {
		return 0, read, fmt.Errorf("failed to read long form TLV length: %w", err)
	}
	read += len(length)

	return int(new(big.Int).SetBytes(length).Int64()), read, nil
}

// Inspect returns human readable information about length prefixer.
func (p *berTLVPrefixer) Inspect() string {
	return "BerTLV"
}

// clearMSB clears the most significant bit at pos in 8. We shift set bit 7
// times (so it becomes 0b10000000). We then flip every bit in the mask with
// the ^ operator (so 0b10000000 becomes 0b01111111). Finally, we use a bitwise
// AND, which doesn't touch the numbers AND'ed with 1, but which will unset the
// value in the mask which is set to 0.
func clearMSB(n uint8) uint8 {
	return n &^ (1 << 7)
}

// setMSB sets the most significant bit of n.
func setMSB(n uint8) uint8 {
	return n | (1 << 7)
}
