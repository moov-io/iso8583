package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVMLHeader(t *testing.T) {
	// t.Run("Pack returns binary encoded length", func(t *testing.T) {
	// 	header := NewBinary2BytesHeader()

	// 	header.SetLength(319)
	// 	var buf bytes.Buffer
	// 	n, err := header.WriteTo(&buf)

	// 	require.NoError(t, err)
	// 	require.Equal(t, 2, n)

	// 	// len 319 encoded in bytes
	// 	require.Equal(t, []byte{0x01, 0x3F}, buf.Bytes())
	// })

	t.Run("Read reads 4 bytes and decode length with session control", func(t *testing.T) {
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

	t.Run("ReadFrom returns error when message length exceeds max message lenght", func(t *testing.T) {
		// header := NewVMLHeader()
	})

	// test
	// 2048 - max message length

}
