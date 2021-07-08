package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCD(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		res, read, err := BCD.Decode([]byte{0x12, 0x34}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("1234"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x01, 0x23}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("123"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x12, 0x30}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("230"), res)
		require.Equal(t, 2, read)
	})

	t.Run("ReadFrom", func(t *testing.T) {
		res, read, err := BCD.DecodeFrom(bytes.NewReader([]byte{0x12, 0x34}), 4)

		require.NoError(t, err)
		require.Equal(t, []byte("1234"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.DecodeFrom(bytes.NewReader([]byte{0x01, 0x23}), 3)

		require.NoError(t, err)
		require.Equal(t, []byte("123"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.DecodeFrom(bytes.NewReader([]byte{0x12, 0x30}), 3)

		require.NoError(t, err)
		require.Equal(t, []byte("230"), res)
		require.Equal(t, 2, read)
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := BCD.Encode([]byte("0110"))

		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x10}, res)

		// right justified by default
		res, err = BCD.Encode([]byte("123"))

		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, res)
	})
}
