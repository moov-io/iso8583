package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHex(t *testing.T) {
	pref := hexFixedPrefixer{}

	dataLen, read, err := pref.DecodeLength(16, []byte("whatever"))

	require.NoError(t, err)
	require.Equal(t, 16, dataLen)
	require.Equal(t, 0, read)
}
