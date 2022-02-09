package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEBCDIC(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		res, read, err := EBCDIC.Decode([]byte{0x12, 0x34}, 2)

		require.NoError(t, err)
		require.Equal(t, []byte{0x12, 0x94}, res)
		require.Equal(t, 2, read)

		_, _, err = EBCDIC.Decode([]byte{0x12, 0x34}, 3)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 2")

		_, _, err = EBCDIC.Decode(nil, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 6, got 0")

	})

	t.Run("Encode", func(t *testing.T) {
		res, err := EBCDIC.Encode([]byte{0x12, 0x94})

		require.NoError(t, err)
		require.Equal(t, []byte{0x12, 0x34}, res)

	})
}
