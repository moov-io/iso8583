package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinary2BytesHeader(t *testing.T) {
	t.Run("Pack returns binary encoded length", func(t *testing.T) {
		header := NewBinary2BytesHeader()

		header.SetLength(319)
		var buf bytes.Buffer
		n, err := header.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, 2, n)

		// len 319 encoded in bytes
		require.Equal(t, []byte{0x01, 0x3F}, buf.Bytes())
	})

	t.Run("Read reads 2 bytes and decode length", func(t *testing.T) {
		header := NewBinary2BytesHeader()

		// len 319 encoded in BCD
		packed := []byte{0x01, 0x3F}
		_, err := header.ReadFrom(bytes.NewReader(packed))

		require.NoError(t, err)
		require.Equal(t, 319, header.Length())
	})
}
