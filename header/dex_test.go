package header

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVisaDEXHeader(t *testing.T) {
	t.Run("Pack returns BCD encoded length", func(t *testing.T) {
		header := NewVisaDEXHeader()

		header.SetLength(115)
		packed, err := header.Pack()

		require.NoError(t, err)
		// len 115 encoded in BCD
		require.Equal(t, []byte{0x01, 0x15}, packed)
	})

	t.Run("Read reads 2 bytes and decode length from BCD", func(t *testing.T) {
		header := NewVisaDEXHeader()

		// len 115 encoded in BCD
		got := []byte{0x01, 0x15}
		_, err := header.Read(bytes.NewReader(got))

		require.NoError(t, err)
		require.Equal(t, 115, header.Length())
	})
}
