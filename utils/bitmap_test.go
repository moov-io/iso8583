package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitmap(t *testing.T) {
	bitmap := NewBitmap(50)
	bitmap.Set(10)
	require.True(t, bitmap.IsSet(10))
	require.False(t, bitmap.IsSet(11))

	src := []byte{232, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	bitmap2 := NewBitmapFromData(src)
	require.Equal(t, src, bitmap2.Bytes())
	require.True(t, bitmap2.IsSet(1))
	require.True(t, bitmap2.IsSet(2))
	require.True(t, bitmap2.IsSet(3))
	require.True(t, bitmap2.IsSet(5))

	bitmap3 := NewBitmap(16)
	bitmap3.Set(1)
	bitmap3.Set(2)
	bitmap3.Set(6)
	bitmap3.Set(7)
	bitmap3.Set(10)
	bitmap3.Set(11)
	bitmap3.Set(15)
	require.Equal(t, []byte{0xC6, 0x62}, bitmap3.Bytes())
}
