package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCD(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		res, err := BCD.Decode([]byte{0x12, 0x34}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("1234"), res)

		res, err = BCD.Decode([]byte{0x01, 0x23}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("123"), res)

		res, err = BCD.Decode([]byte{0x12, 0x30}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("230"), res)
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
