package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yerden/go-util/bcd"
)

func TestTrack2BCD(t *testing.T) {
	t.Run("Decode - sample even length (36 digits -> 18 bytes)", func(t *testing.T) {
		// "4000340000000504D2225111123400001230" (36 nibbles)
		src := []byte{
			0x40, 0x00, 0x34, 0x00, 0x00, 0x00, 0x05, 0x04,
			0xD2, 0x22, 0x51, 0x11, 0x12, 0x34, 0x00, 0x00,
			0x12, 0x30,
		}
		out, read, err := Track2BCD.Decode(src, 36)

		require.NoError(t, err)
		require.Equal(t, 18, read)
		require.Equal(t, []byte("4000340000000504D2225111123400001230"), out)
	})

	t.Run("Decode - odd length with leading zero alignment", func(t *testing.T) {
		// bytes 0x01,0x23 represent “0123” in BCD.
		// length=3 => remove the extra leading nibble => “123”
		out, read, err := Track2BCD.Decode([]byte{0x01, 0x23}, 3)

		require.NoError(t, err)
		require.Equal(t, 2, read)
		require.Equal(t, []byte("123"), out)
	})

	t.Run("Decode - accepts 'D' separator nibble (0xD)", func(t *testing.T) {
		// 0x01,0x2D => "012D"; length=3 => "12D"
		out, read, err := Track2BCD.Decode([]byte{0x01, 0x2D}, 3)

		require.NoError(t, err)
		require.Equal(t, 2, read)
		require.Equal(t, []byte("12D"), out)
	})

	t.Run("Decode - not enough data", func(t *testing.T) {
		// length=4 => requires 2 bytes; only 1 byte provided
		_, _, err := Track2BCD.Decode([]byte{0x12}, 4)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 2, got 1")
	})

	t.Run("Decode - invalid nibble (A..F except D)", func(t *testing.T) {
		// 0x1A contains unmapped nibble 0xA => ErrBadBCD
		_, _, err := Track2BCD.Decode([]byte{0x1A, 0x23}, 4)
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD decoding")
		require.ErrorIs(t, err, bcd.ErrBadBCD)
	})

	t.Run("Decode - ignores trailing bytes beyond ceil(n/2)", func(t *testing.T) {
		// length=4 => reads 2 bytes; the third byte (0xFF) is ignored
		out, read, err := Track2BCD.Decode([]byte{0x21, 0x43, 0xFF}, 4)

		require.NoError(t, err)
		require.Equal(t, 2, read)
		require.Equal(t, []byte("2143"), out)
	})

	t.Run("Encode - round-trip with sample", func(t *testing.T) {
		in := []byte("4000340000000504D2225111123400001230") // 36 digits
		packed, err := Track2BCD.Encode(in)
		require.NoError(t, err)
		require.Len(t, packed, 18) // ceil(36/2)

		out, read, err := Track2BCD.Decode(packed, len(in))
		require.NoError(t, err)
		require.Equal(t, 18, read)
		require.Equal(t, in, out)
	})

	t.Run("Encode - odd length uses leading zero (right-justified)", func(t *testing.T) {
		packed, err := Track2BCD.Encode([]byte("123")) // 3 digits -> 0x01,0x23
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, packed)
	})

	t.Run("Encode - includes 'D' nibble", func(t *testing.T) {
		packed, err := Track2BCD.Encode([]byte("12D3")) // -> 0x12,0xD3
		require.NoError(t, err)
		require.Equal(t, []byte{0x12, 0xD3}, packed)
	})

	t.Run("Encode - invalid chars", func(t *testing.T) {
		_, err := Track2BCD.Encode([]byte("abc")) // 'a', 'b', 'c' are not expected
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD encoding")
		require.ErrorIs(t, err, bcd.ErrBadInput)
	})
}

func FuzzDecodeTrack2BCD(f *testing.F) {
	enc := Track2BCD

	f.Fuzz(func(t *testing.T, data []byte, length int) {
		enc.Decode(data, length)
	})
}
