package encoding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/moov-io/iso8583/utils"
)

// HEX to ASCII encoder
var (
	_ Encoder = (*hexToASCIIEncoder)(nil)
	// BytesToASCIIHex is an encoder that converts bytes into their ASCII
	// representation.  On success, the ASCII representation bytes are returned
	// Don't use this encoder with String, Numeric or Binary fields as packing and
	// unpacking in these fields uses length of value/bytes, so only Pack will be
	// able to write the value correctly.
	BytesToASCIIHex = &hexToASCIIEncoder{}
)

type hexToASCIIEncoder struct{}

// Encode converts bytes into their ASCII representation.  On success, the
// ASCII representation bytes are returned e.g. []byte{0x5F, 0x2A} would be
// converted to []byte("5F2A")
func (e hexToASCIIEncoder) Encode(data []byte) ([]byte, error) {
	out := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(out, data)

	str := string(out)
	str = strings.ToUpper(str)

	return []byte(str), nil
}

// Decodes ASCII hex and returns bytes
// length is number of HEX-digits (two ASCII characters is one HEX digit)
// e.g. []byte("AABBCC") would be converted into []byte{0xAA, 0xBB, 0xCC}
func (e hexToASCIIEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	// to read 8 HEX digits we have to read 16 ASCII chars (bytes)
	read := hex.EncodedLen(length)
	if read > len(data) {
		return nil, 0, errors.New("not enough data to read")
	}

	out := make([]byte, length)

	_, err := hex.Decode(out, data[:read])
	if err != nil {
		return nil, 0, utils.NewSafeError(err, "failed to perform hex decoding")
	}

	return out, read, nil
}

// ASCII To HEX encoder
var (
	_ Encoder = (*asciiToHexEncoder)(nil)
	// ASCIIHexToBytes is an encoder that converts ASCII Hex-digits into a byte slice
	// This encoder is used in TagSpec, BerTLVTag and similar.
	// It shouldn't be used with String, Numeric or Binary fields as packing and unpacking
	// in these fields uses length of value/bytes, so only Unpack will be able to read
	// the value correctly.
	// If you are looking for a way to work with HEX strings, use Hex field instead.
	ASCIIHexToBytes = &asciiToHexEncoder{}
)

type asciiToHexEncoder struct{}

// Encode converts ASCII Hex-digits into a byte slice e.g. []byte("AABBCC")
// would be converted into []byte{0xAA, 0xBB, 0xCC}
func (e asciiToHexEncoder) Encode(data []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(data)))

	_, err := hex.Decode(out, data)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to perform hex decoding")
	}

	return out, nil
}

// Decode converts bytes into their ASCII representation.
// Length is number of ASCII characters (two ASCII characters is one HEX digit)
// On success, the ASCII representation bytes are returned e.g. []byte{0x5F,
// 0x2A} would be converted to []byte("5F2A")
func (e asciiToHexEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	if length > len(data) {
		return nil, 0, errors.New("not enough data to read")
	}

	out := make([]byte, hex.EncodedLen(length))
	hex.Encode(out, data[:length])

	return []byte(strings.ToUpper(string(out))), length, nil
}
