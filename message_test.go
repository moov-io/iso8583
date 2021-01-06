package iso8583

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/field"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	message := NewMessage(Spec87)

	message.Field(0, "0100")
	message.Field(2, "424242424242")

	require.Equal(t, "0100", message.GetString(0))
	require.Equal(t, "424242424242", message.GetString(2))
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

	t.Run("Test packing with typed fields", func(t *testing.T) {
		type ISO87Data struct {
			F2 *field.StringField
			F3 *field.NumericField
			F4 *field.StringField
		}

		message := NewMessage(Spec87)
		err := message.SetData(&ISO87Data{
			F2: field.NewStringValue("4242424242424242"),
			F3: field.NewNumericValue(123456),
			F4: field.NewStringValue("100"),
		})
		require.NoError(t, err)

		fmt.Println(message.Fields)

		message.MTI("0100")

		rawMsg, err := message.Pack()

		require.NoError(t, err)

		wantMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		require.Equal(t, wantMsg, rawMsg)
	})
}
