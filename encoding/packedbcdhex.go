package encoding

import (
	"fmt"

	"github.com/yerden/go-util/bcd"
)

var (
	_ Encoder = (*packedBCDHexEncoder)(nil)

	// PackedBCDHex is a packed (unsigned) BCD encoder that decodes
	// non-decimal nibbles (A-F) into their hexadecimal characters,
	// instead of rejecting them the way the decimal-only BCD encoder
	// does.
	//
	// Like BCD, the length unit of Decode is the number of DIGITS: the
	// packed data occupies ceil(length/2) bytes, and values with an odd
	// number of digits are right-aligned in their packed bytes (the
	// first nibble is fill). Unlike BCD, every nibble 0x0-0xF decodes to
	// its hex character 0-9A-F.
	//
	// It exists for the packed-BCD variable-length fields of payment
	// network specifications that legitimately carry nibbles outside
	// 0-9, the canonical case being VisaNet Track 2 Data (DE35):
	//
	//	Field 35 - Variable length 1 byte, binary + 37 N,
	//	             4-bit BCD (unsigned packed); maximum 20 bytes
	//
	// There the LL prefix counts track 2 DIGITS (up to 37) while the
	// packed data (19 bytes) carries the 'D' separator nibble and the
	// 'F' padding nibbles. BCD errors on those nibbles, and
	// ASCIIHexToBytes consumes LL BYTES (2x the data), so neither can
	// unpack the field. PackedBCDHex together with prefix.Binary.L
	// expresses the field exactly:
	//
	//	&field.Spec{
	//		Length:      37,
	//		Description: "Track 2 Data",
	//		Enc:         encoding.PackedBCDHex,
	//		Pref:        prefix.Binary.L,
	//	}
	//
	// Encode accepts the characters 0-9, A-F and a-f; a value with an
	// odd number of characters is left-padded with a fill nibble, as
	// BCD does. Values containing track separators other than the
	// nibble-encoded 'D' (such as '=') cannot be packed and are
	// rejected. Errors never include the input data.
	PackedBCDHex = &packedBCDHexEncoder{}
)

// packedBCDHexDigits maps a nibble to its uppercase hexadecimal
// character, matching the output casing of ASCIIHexToBytes.
const packedBCDHexDigits = "0123456789ABCDEF"

type packedBCDHexEncoder struct{}

// Encode packs the hexadecimal characters of src into unsigned packed
// BCD, two characters per byte. Odd lengths are left-padded with a fill
// nibble, mirroring the right-alignment of bcdEncoder.Encode. The
// returned error identifies only the position of an offending
// character in src, never its value.
func (e *packedBCDHexEncoder) Encode(src []byte) ([]byte, error) {
	// left-pad odd lengths with a fill character, as BCD does
	padded := src
	if len(src)%2 != 0 {
		padded = append([]byte("0"), src...)
	}

	packed := make([]byte, len(padded)/2)
	for i, c := range padded {
		var nibble byte
		switch {
		case c >= '0' && c <= '9':
			nibble = c - '0'
		case c >= 'A' && c <= 'F':
			nibble = c - 'A' + 10
		case c >= 'a' && c <= 'f':
			nibble = c - 'a' + 10
		default:
			// do not echo c: values such as track data are sensitive
			pos := i
			if len(src)%2 != 0 {
				pos-- // compensate for the prepended fill character
			}
			return nil, fmt.Errorf("failed to perform packed BCD-HEX encoding: expected a hexadecimal character at position %d", pos)
		}

		if i%2 == 0 {
			packed[i/2] = nibble << 4
		} else {
			packed[i/2] |= nibble
		}
	}

	return packed, nil
}

// Decode unpacks length DIGITS from ceil(length/2) bytes, mapping every
// nibble to its hexadecimal character. Odd lengths are right-aligned
// exactly like bcdEncoder.Decode: the leading fill nibble is skipped.
func (e *packedBCDHexEncoder) Decode(src []byte, length int) ([]byte, int, error) {
	// length should be positive
	if length < 0 {
		return nil, 0, fmt.Errorf("length should be positive, got %d", length)
	}

	// as with BCD the decoded digit count is even
	decodedLen := length
	if length%2 != 0 {
		decodedLen += 1
	}

	// how many bytes we will read
	read := bcd.EncodedLen(decodedLen)

	if len(src) < read {
		return nil, 0, fmt.Errorf("not enough data to decode. expected len %d, got %d", read, len(src))
	}

	dst := make([]byte, decodedLen)
	for i, b := range src[:read] {
		dst[2*i] = packedBCDHexDigits[b>>4]
		dst[2*i+1] = packedBCDHexDigits[b&0x0f]
	}

	// because packed BCD is right aligned, we skip first bytes and
	// read only what we need
	// e.g. 0643 => 643
	return dst[decodedLen-length:], read, nil
}
