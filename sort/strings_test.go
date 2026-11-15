package sort

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortStrings(t *testing.T) {
	x := []string{"1", "2", "11"}
	Strings(x)
	require.Equal(t, []string{"1", "11", "2"}, x)
}

func TestSortStringsByInt(t *testing.T) {
	x := []string{"11", "5", "1"}
	StringsByInt(x)
	require.Equal(t, []string{"1", "5", "11"}, x)
}

func TestSortStringsByHex(t *testing.T) {
	t.Run("sorts by numeric hex value", func(t *testing.T) {
		x := []string{"B0", "10", "ABCD"}
		StringsByHex(x)
		require.Equal(t, []string{"10", "B0", "ABCD"}, x)
	})

	t.Run("ignores leading zeros and hex case", func(t *testing.T) {
		x := []string{"0010", "0001", "ab"}
		StringsByHex(x)
		require.Equal(t, []string{"0001", "0010", "ab"}, x)
	})

	t.Run("falls back to string comparison for invalid hex", func(t *testing.T) {
		x := []string{"xyz", "10", "ABC"}
		StringsByHex(x)
		require.Equal(t, []string{"10", "ABC", "xyz"}, x)
	})
}
