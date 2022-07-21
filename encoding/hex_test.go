package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexToASCIIEncoder(t *testing.T) {
	enc := BytesToASCIIHex

	got, read, err := enc.Decode([]byte("aabbcc"), 3)
	require.NoError(t, err)
	require.Equal(t, 6, read)
	require.Equal(t, []byte{0xAA, 0xBB, 0xCC}, got)

	_, _, err = enc.Decode(nil, 3)
	require.Error(t, err)
	require.EqualError(t, err, "not enough data to read")

	_, _, err = enc.Decode([]byte("nothex"), 3)
	require.Error(t, err)
	require.EqualError(t, err, "failed to perform hex decoding")

	got, err = enc.Encode([]byte{0xAA, 0xBB, 0xCC})
	require.NoError(t, err)
	require.Equal(t, []byte("AABBCC"), got)
}

func TestASCIIToHexEncoder(t *testing.T) {
	enc := ASCIIHexToBytes

	got, read, err := enc.Decode([]byte{0xAA, 0xBB, 0xCC}, 3)
	require.NoError(t, err)
	require.Equal(t, []byte("AABBCC"), got)
	require.Equal(t, 3, read)

	got, err = enc.Encode([]byte("aabbcc"))
	require.NoError(t, err)
	require.Equal(t, []byte{0xAA, 0xBB, 0xCC}, got)

	_, err = enc.Encode([]byte("nothex"))
	require.Error(t, err)
	require.EqualError(t, err, "failed to perform hex decoding")
}
