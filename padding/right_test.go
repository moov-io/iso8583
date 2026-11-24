package padding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRightPadder(t *testing.T) {
	padder := NewRightPadder('0')

	t.Run("Pad", func(t *testing.T) {
		str := []byte("12345")
		want := []byte("1234500000")

		got := padder.Pad(str, 10)

		require.Equal(t, want, got)
	})

	t.Run("Pad preserves input longer than the target", func(t *testing.T) {
		require.Equal(t, []byte("12345"), padder.Pad([]byte("12345"), 3))
	})

	t.Run("Inspect returns the padding rune as bytes", func(t *testing.T) {
		require.Equal(t, []byte("é"), Right('é').Inspect())
	})

	t.Run("Unpad", func(t *testing.T) {
		str := []byte("1234500000")
		want := []byte("12345")

		got := padder.Unpad(str)

		require.Equal(t, want, got)
	})
}
