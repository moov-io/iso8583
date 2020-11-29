package iso8583

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitmap(t *testing.T) {
	bitmap := NewBitmap(50)
	bitmap.Set(10)
	require.True(t, bitmap.IsSet(10))
	require.False(t, bitmap.IsSet(11))
}
