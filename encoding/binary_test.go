package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinary(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03, 0x04}

		out, read, err := Binary.Decode(data, 2)
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x02}, out)
		require.Equal(t, 2, read)

		// the result must not alias the input
		out[0] = 0xff
		require.Equal(t, byte(0x01), data[0])

		out, read, err = Binary.Decode(data, 0)
		require.NoError(t, err)
		require.Nil(t, out)
		require.Equal(t, 0, read)

		_, _, err = Binary.Decode(data, 5)
		require.EqualError(t, err, "failed to perform binary decoding: length 5 exceeds the data size 4")
	})

	t.Run("Encode", func(t *testing.T) {
		data := []byte{0x01, 0x02}

		out, err := Binary.Encode(data)
		require.NoError(t, err)
		require.Equal(t, data, out)

		out[0] = 0xff
		require.Equal(t, byte(0x01), data[0])
	})
}

func FuzzDecodeBinary(f *testing.F) {
	enc := &binaryEncoder{}

	f.Fuzz(func(t *testing.T, data []byte, length int) {
		enc.Decode(data, length)
	})
}
