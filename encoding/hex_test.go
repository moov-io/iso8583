package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHex(t *testing.T) {
	enc := &hexEncoder{}

	got, err := enc.Decode([]byte("aabbcc"), 0)
	require.NoError(t, err)
	require.Equal(t, []byte{0xAA, 0xBB, 0xCC}, got)

	got, err = enc.Encode([]byte{0xAA, 0xBB, 0xCC})
	require.NoError(t, err)
	require.Equal(t, []byte("aabbcc"), got)
}
