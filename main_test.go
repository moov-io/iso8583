package iso8583

import (
	"testing"
)

func TestISO8583(t *testing.T) {
	// specTest := &spec.MessageSpec{
	// 	Fields: map[int]spec.Packer{
	// 		0: spec.NewField(4, "Message Type Indicator", encoding.ASCII, prefix.ASCII.Fixed),

	// 		// Bitmap, 16 bytes, fixed
	// 		1: spec.Bitmap(16, "Bitmap", encoding.Hex, prefix.Hex.Fixed),

	// 		// LLVAR19
	// 		2: &spec.Field{
	// 			Length:      19,
	// 			Description: "Primary Account Number",
	// 			Enc:         encoding.ASCII,
	// 			Pref:        prefix.ASCII.LL,
	// 		},

	// 		// 6 bytes, fixed
	// 		3: &spec.Field{
	// 			Length:      6,
	// 			Description: "Processing Code",
	// 			Enc:         encoding.ASCII,
	// 			Pref:        prefix.ASCII.Fixed,
	// 			Pad:         padding.Left('0'),
	// 		},

	// 		// 12 bytes, fixed
	// 		4: &spec.Field{
	// 			Length:      12,
	// 			Description: "Transaction Amount",
	// 			Enc:         encoding.ASCII,
	// 			Pref:        prefix.ASCII.Fixed,
	// 			Pad:         padding.Left('0'),
	// 		},

	// 		5: &spec.Field{},
	// 	},
	// }

	// message := NewMessage(specTest)
	// message.MTI("0100")
	// message.Field(2, "4242424242424242")
	// message.Field(3, "123456")
	// message.Field(4, "100")

	// got, err := message.Pack()

	// want := "01007000000000000000164242424242424242123456000000000100"
	// require.NoError(t, err)
	// require.NotNil(t, got)
	// require.Equal(t, want, string(got))

	// message = NewMessage(specTest)
	// message.Unpack([]byte(want))

	// require.Equal(t, "0100", message.GetMTI())
	// require.Equal(t, "4242424242424242", message.GetString(2))
	// require.Equal(t, "123456", message.GetString(3))
	// require.Equal(t, "100", message.GetString(4))
}
