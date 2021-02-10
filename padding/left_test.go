package padding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeftPadder(t *testing.T) {
	padder := NewLeftPadder('0')

	t.Run("Pad", func(t *testing.T) {
		str := []byte("12345")
		want := []byte("0000012345")

		got := padder.Pad(str, 10)

		require.Equal(t, want, got)
	})

	t.Run("Unpad", func(t *testing.T) {
		str := []byte("0000012345")
		want := []byte("12345")

		got := padder.Unpad(str)

		require.Equal(t, want, got)
	})
}
