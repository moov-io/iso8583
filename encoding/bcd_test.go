package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yerden/go-util/bcd"
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

		res, read, err = BCD.Decode([]byte{0x21, 0x43, 0x55}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("2143"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x21, 0x43, 0xff}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("2143"), res)
		require.Equal(t, 2, read)

		_, _, err = BCD.Decode([]byte{0x21, 0x43}, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 2")

		_, _, err = BCD.Decode(nil, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 0")

		_, _, err = BCD.Decode([]byte{0xAB, 0xCD}, 4)
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD decoding")
		require.ErrorIs(t, err, bcd.ErrBadBCD)
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := BCD.Encode([]byte("0110"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x10}, res)

		// right justified by default
		res, err = BCD.Encode([]byte("123"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, res)

		_, err = BCD.Encode([]byte("abc"))
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD encoding")
		require.ErrorIs(t, err, bcd.ErrBadInput)
	})
}
