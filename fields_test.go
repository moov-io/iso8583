package iso8583

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestField(t *testing.T) {
	f := NewField(1, []byte("hello"))

	require.Equal(t, "hello", f.String())
	require.Equal(t, []byte("hello"), f.Bytes())
}
