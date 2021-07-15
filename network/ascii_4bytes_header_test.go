package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestASCII4BytesHeader(t *testing.T) {
	t.Run("Pack returns ASCII encoded length", func(t *testing.T) {
		header := NewASCII4BytesHeader()

		header.SetLength(115)
		var buf bytes.Buffer
		n, err := header.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, 4, n)
		// len 115 encoded in ASCII
		require.Equal(t, "0115", buf.String())
	})

	t.Run("Read reads 4 bytes and decode length from ASCII", func(t *testing.T) {
		header := NewASCII4BytesHeader()

		// len 115 encoded in ASCII
		packed := []byte("0115")
		_, err := header.ReadFrom(bytes.NewReader(packed))

		require.NoError(t, err)
		require.Equal(t, 115, header.Length())
	})

	t.Run("Read returns error when not enough data to read", func(t *testing.T) {
		header := NewASCII4BytesHeader()

		// len 115 encoded in ASCII
		packed := []byte("011")
		_, err := header.ReadFrom(bytes.NewReader(packed))

		require.Error(t, err)
	})
}
