package prefixer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHex(t *testing.T) {
	pref := hexFixedPrefixer{
		Len: 16,
	}

	dataLen, err := pref.DecodeLength([]byte("whatever"))

	require.NoError(t, err)
	require.Equal(t, 32, dataLen)

}
