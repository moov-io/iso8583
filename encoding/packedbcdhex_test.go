package encoding

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackedBCDHex(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		// decimal nibbles decode exactly like BCD
		res, read, err := PackedBCDHex.Decode([]byte{0x12, 0x34}, 4)
		require.NoError(t, err)
		require.Equal(t, []byte("1234"), res)
		require.Equal(t, 2, read)

		// odd digit counts are right-aligned: the leading fill nibble
		// is skipped, the same way 0643 => 643 in BCD
		res, read, err = PackedBCDHex.Decode([]byte{0x01, 0x23}, 3)
		require.NoError(t, err)
		require.Equal(t, []byte("123"), res)
		require.Equal(t, 2, read)

		res, read, err = PackedBCDHex.Decode([]byte{0x12, 0x30}, 3)
		require.NoError(t, err)
		require.Equal(t, []byte("230"), res)
		require.Equal(t, 2, read)

		// non-decimal nibbles A-F decode to their hex characters,
		// where BCD would fail with ErrBadBCD
		res, read, err = PackedBCDHex.Decode([]byte{0xAB, 0xCD}, 4)
		require.NoError(t, err)
		require.Equal(t, []byte("ABCD"), res)
		require.Equal(t, 2, read)

		// 'F' fill nibbles are preserved as the 'F' character
		res, read, err = PackedBCDHex.Decode([]byte{0x21, 0x43, 0xFF}, 4)
		require.NoError(t, err)
		require.Equal(t, []byte("2143"), res)
		require.Equal(t, 2, read)

		// the trailing fill nibble of an odd value decodes too
		res, read, err = PackedBCDHex.Decode([]byte{0x01, 0x2F}, 3)
		require.NoError(t, err)
		require.Equal(t, []byte("12F"), res)
		require.Equal(t, 2, read)

		_, _, err = PackedBCDHex.Decode([]byte{0x21, 0x43}, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 2")

		_, _, err = PackedBCDHex.Decode(nil, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 0")

		_, _, err = PackedBCDHex.Decode([]byte{0x12}, -1)
		require.Error(t, err)
		require.EqualError(t, err, "length should be positive, got -1")
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := PackedBCDHex.Encode([]byte("0110"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x10}, res)

		// right justified by default, as BCD does
		res, err = PackedBCDHex.Encode([]byte("123"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, res)

		// non-decimal digits pack into nibbles A-F
		res, err = PackedBCDHex.Encode([]byte("ABCD"))
		require.NoError(t, err)
		require.Equal(t, []byte{0xAB, 0xCD}, res)

		// lowercase input packs like uppercase
		res, err = PackedBCDHex.Encode([]byte("abcd"))
		require.NoError(t, err)
		require.Equal(t, []byte{0xAB, 0xCD}, res)

		_, err = PackedBCDHex.Encode([]byte("abcg"))
		require.Error(t, err)
		// the error names the position, never the offending character
		require.EqualError(t, err, "failed to perform packed BCD-HEX encoding: expected a hexadecimal character at position 3")

		// the ISO track separator '=' is not a hex nibble and is rejected
		_, err = PackedBCDHex.Encode([]byte("1234=678"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected a hexadecimal character at position 4")
	})

	t.Run("round trip", func(t *testing.T) {
		for _, s := range []string{
			"0", "12", "0F", "2D", "ABCDEF", "0123456789",
			strings.Repeat("0123456789ABCDEF", 4),
		} {
			packed, err := PackedBCDHex.Encode([]byte(s))
			require.NoError(t, err)

			unpacked, read, err := PackedBCDHex.Decode(packed, len(s))
			require.NoError(t, err)
			require.Equal(t, len(packed), read)
			require.Equal(t, s, string(unpacked))
		}
	})
}
