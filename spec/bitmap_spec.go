package spec

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/stretchr/testify/require"
)

func TestBitmapField(t *testing.T) {
	field := &bitmapField{"Bitmap", encoding.Hex, prefixer.Hex.Fixed(16)}

	// pack
	got, err := field.Pack([]byte{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF, 0xAB})
	want := []byte("abcdefabcdefabcdefabcdefabcdefab")
	require.NoError(t, err)
	require.Equal(t, want, got)

	// unpack
	// when only primari bitmap presents
	// we should read only first 8 bytes
	got, length, err := field.Unpack([]byte("68000000000000000000000000000000123456"))
	want = []byte{104, 0, 0, 0, 0, 0, 0, 0}
	require.Equal(t, 16, length)
	require.Len(t, got, 8)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// when secondary primari bitmap presents
	// we should read 16 bytes
	got, length, err = field.Unpack([]byte("E8000000000000000000000000000000aa"))
	want = []byte{232, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	require.Equal(t, 32, length)
	require.Len(t, got, 16)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
