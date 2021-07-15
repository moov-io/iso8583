package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCD2BytesHeader(t *testing.T) {
	t.Run("Pack returns BCD encoded length", func(t *testing.T) {
		header := NewBCD2BytesHeader()

		header.SetLength(115)
		var buf bytes.Buffer
		n, err := header.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, 2, n)
		// len 115 encoded in BCD
		require.Equal(t, []byte{0x01, 0x15}, buf.Bytes())
	})

	t.Run("Read reads 2 bytes and decode length from BCD", func(t *testing.T) {
		header := NewBCD2BytesHeader()

		// len 115 encoded in BCD
		packed := []byte{0x01, 0x15}
		_, err := header.ReadFrom(bytes.NewReader(packed))

		require.NoError(t, err)
		require.Equal(t, 115, header.Length())
	})

	t.Run("Read returns error when not enough data to read", func(t *testing.T) {
		header := NewBCD2BytesHeader()

		// len 115 encoded in BCD
		packed := []byte{0x01}
		_, err := header.ReadFrom(bytes.NewReader(packed))

		require.Error(t, err)
	})
}
