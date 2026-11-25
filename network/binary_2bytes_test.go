package network

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type errorWriter struct{ err error }

func (w errorWriter) Write([]byte) (int, error) { return 0, w.err }

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

	t.Run("Read rejects a truncated header", func(t *testing.T) {
		header := NewBinary2BytesHeader()
		n, err := header.ReadFrom(bytes.NewReader([]byte{0x01}))
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		require.Zero(t, n)
	})

	t.Run("Write preserves the writer error", func(t *testing.T) {
		header := NewBinary2BytesHeader()
		require.NoError(t, header.SetLength(1))
		cause := errors.New("test writer error")
		n, err := header.WriteTo(errorWriter{err: cause})
		require.ErrorIs(t, err, cause)
		require.Zero(t, n)
	})

	t.Run("maximum length and overflow", func(t *testing.T) {
		header := NewBinary2BytesHeader()
		require.NoError(t, header.SetLength(math.MaxUint16))
		var buf bytes.Buffer
		n, err := header.WriteTo(&buf)
		require.NoError(t, err)
		require.Equal(t, 2, n)
		require.Equal(t, []byte{0xFF, 0xFF}, buf.Bytes())
		require.ErrorContains(t, header.SetLength(math.MaxUint16+1), "exceeds max length")
		require.Equal(t, math.MaxUint16, header.Length())
	})
}
