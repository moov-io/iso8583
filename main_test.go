package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/moov-io/iso8583/spec"
	"github.com/stretchr/testify/require"
)

// Encodings:
// ASCII
// ASCIIHEX
// Binary
// BCD

// Prefixers:
// ASCII(x) x - 1, 2, 3
// BCD(x) x - 1, 2, 3

// Padders:
// Left('0')
// Right('0')
// None()

func TestISO8583(t *testing.T) {
	specTest := &spec.MessageSpec{
		Fields: map[int]spec.Packer{
			// 0: spec.MTI(4, "Message Type Indicator", encoding.ASCII),
			// 1: spec.Bitmap(16, "Bitmap", encoding.ASCII),

			// LLVAR19
			2: spec.NewField("Primary Account Number", encoding.ASCII, prefixer.ASCII.LL(19)),

			// Processing Code, Numeric, 6 bytes, fixed
			// 3: spec.NewField(6, "123456", encoding.ASCII, prefixer.None, padder.None, validator.Numeric),

			// Settlement Amount, Numeric, 12 bytes, fixed
			// 5: spec.NewField(12, "Settlement Amount", encoding.ASCII, prefixer.None, padder.Left('0'), validator.Numeric),
		},
	}

	message := NewMessage(specTest)
	message.Field(0, "0100")
	message.Field(2, "4242424242424242")
	message.Field(3, "123456")
	message.Field(5, "100")

	got, err := message.Pack()
	require.NoError(t, err)
	require.NotNil(t, raw)

	want := []byte{
		// 0: MTI
		48, 49, 50, 51,
		// 1: Bitmap
		104, 0, 0, 0, 0, 0, 0, 0,
		// 2: Account Number - 4242424242424242
		49, 54, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50,
		// 3: Processing Code - 123456
		49, 50, 51, 52, 53, 54,
		// 5: Settlement Amount - 100, left padded with 0
		48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 48, 48,
	}

	require.Equal(t, want, got)

	// message = NewMessage(specTest)
	// message.Unpack(raw)

	// require.Equal(t, "0100", message.GetString(0))
	// require.Equal(t, "424242424242", message.GetString(2))
}
