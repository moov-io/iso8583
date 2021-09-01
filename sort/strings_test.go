package sort

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringsByInt(t *testing.T) {
	x := []string{"11", "5", "1"}
	StringsByInt(x)
	require.Equal(t, []string{"1", "5", "11"}, x)
}

func TestStringsByHex(t *testing.T) {
	x := []string{"B0", "10", "ABCD"}
	StringsByHex(x)
	require.Equal(t, []string{"10", "B0", "ABCD"}, x)
}
