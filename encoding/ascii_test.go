package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestASCII(t *testing.T) {
	enc := &asciiEncoder{}

	t.Run("Decode", func(t *testing.T) {
		res, read, err := enc.Decode([]byte("hello"), 5)

		require.NoError(t, err)
		require.Equal(t, []byte("hello"), res)
		require.Equal(t, 5, read)

		_, _, err = enc.Decode([]byte("hello, 世界!"), 10)
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform ASCII decoding")

		_, _, err = enc.Decode([]byte("hello"), 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 6, got 5")

		_, _, err = enc.Decode(nil, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 6, got 0")
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := enc.Encode([]byte("hello"))

		require.NoError(t, err)
		require.Equal(t, []byte("hello"), res)

		_, err = enc.Encode([]byte("hello, 世界!"))
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform ASCII encoding")
	})
}
