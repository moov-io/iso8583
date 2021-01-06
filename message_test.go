package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	spec := &MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewStringField(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmapField(&field.Spec{
				Length:      16,
				Description: "Bitmap",
				Enc:         encoding.Hex,
				Pref:        prefix.Hex.Fixed,
			}),
			2: field.NewStringField(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			3: field.NewNumericField(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			4: field.NewStringField(&field.Spec{
				Length:      12,
				Description: "Transaction Amount",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}

	t.Run("Test packing and unpacking untyped fields", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		message.Field(2, "4242424242424242")
		message.Field(3, "123456")
		message.Field(4, "100")

		got, err := message.Pack()

		want := "01007000000000000000164242424242424242123456000000000100"
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, string(got))

		message = NewMessage(spec)
		message.Unpack([]byte(want))

		require.Equal(t, "0100", message.GetMTI())
		require.Equal(t, "4242424242424242", message.GetString(2))
		require.Equal(t, "123456", message.GetString(3))
		require.Equal(t, "100", message.GetString(4))
	})

	t.Run("Test unpacking with typed fields", func(t *testing.T) {
		type ISO87Data struct {
			F2 *field.StringField
			F3 *field.NumericField
			F4 *field.StringField
		}

		message := NewMessage(spec)
		message.SetData(&ISO87Data{})

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		require.Equal(t, "4242424242424242", message.GetString(2))
		require.Equal(t, "123456", message.GetString(3))
		require.Equal(t, "100", message.GetString(4))

		data := message.Data().(*ISO87Data)

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

		message := NewMessage(spec)
		message.MTI("0100")
		err := message.SetData(&ISO87Data{
			F2: field.NewStringValue("4242424242424242"),
			F3: field.NewNumericValue(123456),
			F4: field.NewStringValue("100"),
		})
		require.NoError(t, err)

		rawMsg, err := message.Pack()

		require.NoError(t, err)

		wantMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		require.Equal(t, wantMsg, rawMsg)
	})
}
