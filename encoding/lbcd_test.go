package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLBCD(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		res, read, err := LBCD.Decode([]byte{0x12, 0x30}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("1230"), res)
		require.Equal(t, 2, read)

		res, read, err = LBCD.Decode([]byte{0x12, 0x30}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("123"), res)
		require.Equal(t, 2, read)
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := LBCD.Encode([]byte("123"))

		require.NoError(t, err)
		require.Equal(t, []byte{0x12, 0x30}, res)
	})
}
