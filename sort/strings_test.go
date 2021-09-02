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
	x := []string{"B0", "10", "ABCD"}
	StringsByHex(x)
	require.Equal(t, []string{"10", "B0", "ABCD"}, x)
}
