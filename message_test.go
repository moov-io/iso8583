package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/field"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	message := NewMessage(Spec87)

	message.Field(0, "0100")
	message.Field(2, "424242424242")
	message.BinaryField(3, []byte{0x12, 0x34})

	require.Equal(t, "0100", message.GetString(0))
	require.Equal(t, "424242424242", message.GetString(2))
	require.Equal(t, []byte{0x12, 0x34}, message.GetBytes(3))
}

func TestMessageData(t *testing.T) {
	rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")

	t.Run("Test unpacking with typed fields", func(t *testing.T) {
		type ISO87Data struct {
			F2 *field.StringField
			F3 *field.NumericField
			F4 *field.StringField
		}

		message := NewMessage(Spec87)
		message.Data = &ISO87Data{}

		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		require.Equal(t, "4242424242424242", message.GetString(2))
		require.Equal(t, "123456", message.GetString(3))
		require.Equal(t, "100", message.GetString(4))

		data := message.Data.(*ISO87Data)

		require.Equal(t, "4242424242424242", data.F2.Value)
		require.Equal(t, 123456, data.F3.Value)
		require.Equal(t, "100", data.F4.Value)
	})
}
