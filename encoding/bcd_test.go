package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCD(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		res, err := BCD.Decode([]byte{0x12, 0x34})

		require.NoError(t, err)
		require.Equal(t, []byte("1234"), res)
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := BCD.Encode([]byte("0110"))

		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x10}, res)
	})
}
