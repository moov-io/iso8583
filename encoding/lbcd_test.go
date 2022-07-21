package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yerden/go-util/bcd"
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

		_, _, err = LBCD.Decode([]byte{0x12, 0x30}, 5)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 2")

		_, _, err = LBCD.Decode(nil, 5)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 0")

		_, _, err = LBCD.Decode([]byte{0xAB, 0xCD}, 4)
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD decoding")
		require.ErrorIs(t, err, bcd.ErrBadBCD)
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := LBCD.Encode([]byte("123"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x12, 0x30}, res)

		_, err = LBCD.Encode([]byte("abc"))
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD encoding")
		require.ErrorIs(t, err, bcd.ErrBadInput)
	})
}
