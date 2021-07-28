package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVMLHeader(t *testing.T) {
	t.Run("WriteTo writes binary encoded length into writer", func(t *testing.T) {
		header := NewVMLHeader()

		header.SetLength(15)
		var buf bytes.Buffer
		n, err := header.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, 4, n)

		// len 15 encoded in 2 bytes + reserved two bytes 0x00
		require.Equal(t, []byte{0x00, 0x0F, 0x00, 0x00}, buf.Bytes())
	})

	t.Run("WriteTo returns error when message length exceeds max message length", func(t *testing.T) {
		header := NewVMLHeader()
		header.SetLength(MaxMessageLength + 1)

		_, err := header.WriteTo(&bytes.Buffer{})

		require.Error(t, err)
	})

	t.Run("ReadFrom reads 4 bytes and decode length with session control", func(t *testing.T) {
		header := NewVMLHeader()

		// Encoded:
		// * 2 bytes len 15
		// * one reserved byte
		// * final byte:
		//   - left nibble, optional for endpoint generated message, message format indicator: 2 - session control
		//   - right nibble (platform indicator): 0 - direct member
		// 000F0020
		packed := []byte{0x00, 0x0F, 0x00, 0x20}
		read, err := header.ReadFrom(bytes.NewReader(packed))

		require.NoError(t, err)
		require.Equal(t, 15, header.Length())
		require.Equal(t, 4, read)

		require.True(t, header.IsSessionControl)
	})

	t.Run("ReadFrom returns error when message length exceeds max message length", func(t *testing.T) {
		header := NewVMLHeader()
		packed := []byte{0xFF, 0xFF, 0x00, 0x20}

		_, err := header.ReadFrom(bytes.NewReader(packed))

		require.Error(t, err)
	})
}
