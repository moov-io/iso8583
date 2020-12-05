package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/moov-io/iso8583/spec"
	"github.com/stretchr/testify/require"
)

func TestISO8583(t *testing.T) {
	specTest := &spec.MessageSpec{
		Fields: map[int]spec.Packer{
			0: spec.NewField("Message Type Indicator", encoding.ASCII, prefixer.ASCII.Fixed(4)),

			// Bitmap, 16 bytes, fixed
			1: spec.Bitmap("Bitmap", encoding.Hex, prefixer.Hex.Fixed(16)),

			// LLVAR19
			2: spec.NewField("Primary Account Number", encoding.ASCII, prefixer.ASCII.LL(19)),

			// 6 bytes, fixed
			// 3: spec.NewField("Processing Code", encoding.ASCII, prefixer.ASCII.Fixed(6)),
			3: &spec.Field{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefixer.ASCII.Fixed(6),
				Pad:         padding.Left('0'),
			},

			// 12 bytes, fixed
			// 4: spec.NewField("Transaction Amount", encoding.ASCII, prefixer.ASCII.Fixed(12)),
			4: &spec.Field{
				Length:      12,
				Description: "Transaction Amount",
				Enc:         encoding.ASCII,
				Pref:        prefixer.ASCII.Fixed(12),
				Pad:         padding.Left('0'),
			},

			5: &spec.Field{},
		},
	}

	message := NewMessage(specTest)
	message.MTI("0100")
	message.Field(2, "4242424242424242")
	message.Field(3, "123456")
	message.Field(4, "100")

	got, err := message.Pack()

	want := "01007000000000000000164242424242424242123456000000000100"
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, want, string(got))

	message = NewMessage(specTest)
	message.Unpack([]byte(want))

	require.Equal(t, "0100", message.GetMTI())
	require.Equal(t, "4242424242424242", message.GetString(2))
	require.Equal(t, "123456", message.GetString(3))
	require.Equal(t, "100", message.GetString(4))
}
