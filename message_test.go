package iso8583

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	message := NewMessage()

	message.Field(0, "0100")
	message.Field(2, "424242424242")
	message.BinaryField(3, []byte{0x12, 0x34})

	require.Equal(t, "0100", message.GetString(0))
	require.Equal(t, "424242424242", message.GetString(2))
	require.Equal(t, []byte{0x12, 0x34}, message.GetBytes(3))
}
