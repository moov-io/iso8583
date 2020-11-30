package spec

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/stretchr/testify/require"
)

func TestFieldPacker(t *testing.T) {
	t.Run("Bitmap field", func(t *testing.T) {
		field := &bitmapFieldDefinition{"Bitmap", encoding.Hex, prefixer.Hex.Fixed(16)}

		// pack
		got, err := field.Pack([]byte{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB})
		want := []byte("abcdefabcdefabcdefabcdefabcdefab")
		require.NoError(t, err)
		require.Equal(t, want, got)

		// unpack
		// when only primari bitmap presents
		// we read only first 8 bytes
		got, err = field.Unpack([]byte("68000000000000000000000000000000123456"))
		want = []byte{104, 0, 0, 0, 0, 0, 0, 0}
		require.NoError(t, err)
		require.Equal(t, want, got)

		// when secondary primari bitmap presents
		got, err = field.Unpack([]byte("E8000000000000000000000000000000aa"))
		want = []byte{232, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("ASCII VAR field", func(t *testing.T) {
		field := &fieldDefinition{"Primary Account Number", encoding.ASCII, prefixer.ASCII.LL(19)}

		// pack
		got, err := field.Pack([]byte("4242424242424242"))
		want := []byte{49, 54, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50}
		require.NoError(t, err)
		require.Equal(t, want, got)

		// unpack
		got, err = field.Unpack([]byte{49, 54, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50})
		want = []byte("4242424242424242")
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("ASCII Fixed field", func(t *testing.T) {
		field := &fieldDefinition{"Processing code", encoding.ASCII, prefixer.ASCII.Fixed(6)}

		// pack
		got, err := field.Pack([]byte("123456"))
		want := []byte("123456")
		require.NoError(t, err)
		require.Equal(t, want, got)

		// unpack
		got, err = field.Unpack([]byte("123456"))
		want = []byte("123456")
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}
