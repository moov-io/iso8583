package header

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseHeader(t *testing.T) {
	t.Run("Pack returns ASCII encoded length", func(t *testing.T) {
		header := NewBaseHeader()

		header.SetLength(115)
		packed, err := header.Pack()

		require.NoError(t, err)
		// len 115 encoded in BCD
		require.Equal(t, []byte("0115"), packed)
	})

	t.Run("Read reads 4 bytes and decode length from ASCII", func(t *testing.T) {
		header := NewBaseHeader()

		// len 115 encoded in BCD
		packed := []byte("0115")
		_, err := header.Read(bytes.NewReader(packed))

		require.NoError(t, err)
		require.Equal(t, 115, header.Length())
	})

	t.Run("Read returns error when not enough data to read", func(t *testing.T) {
		header := NewBaseHeader()

		// len 115 encoded in BCD
		packed := []byte("011")
		_, err := header.Read(bytes.NewReader(packed))

		require.Error(t, err)
	})
}
