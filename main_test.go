package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
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
			3: spec.NewField("Processing Code", encoding.ASCII, prefixer.ASCII.Fixed(6)),

			// 12 bytes, fixed
			4: spec.NewField("Transaction Amount", encoding.ASCII, prefixer.ASCII.Fixed(12)),
		},
	}

	message := NewMessage(specTest)
	message.Field(0, "0100")
	message.Field(2, "4242424242424242")
	message.Field(3, "123456")
	message.Field(4, "000000000100")

	want := []byte{
		// 0: MTI
		48, 49, 48, 48,
		// 1: Bitmap 01110000, 0000...
		55, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		// 2: Account Number - 4242424242424242
		49, 54, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50,
		// 3: Processing Code - 123456
		49, 50, 51, 52, 53, 54,
		// 4: Settlement Amount - 100, left padded with 0
		48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 48, 48,
	}
	wantStr := "010070000000000000000000000000000000164242424242424242123456000000000100"

	got, err := message.Pack()

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, want, got)
	require.Equal(t, wantStr, string(got))

	// message = NewMessage(specTest)
	// message.Unpack(raw)

	// require.Equal(t, "0100", message.GetString(0))
	// require.Equal(t, "424242424242", message.GetString(2))
}
