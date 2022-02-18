package iso8583

import (
	"bytes"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {
	spec := &MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Description: "Bitmap",
				Enc:         encoding.BytesToASCIIHex,
				Pref:        prefix.Hex.Fixed,
			}),
			2: field.NewString(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
		},
	}

	message := NewMessage(spec)
	message.MTI("0100")
	message.Field(2, "4242424242424242")
	message.Pack() // to generate bitmap

	out := bytes.NewBuffer([]byte{})
	require.NotPanics(t, func() {
		Describe(message, out)
	})

	expectedOutput := `ISO 8583 Message:
MTI...........................: 0100
Bitmap........................: 400000000000000000000000000000000000000000000000
Bitmap bits...................: 01000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
F000 Message Type Indicator...: 0100
F002 Primary Account Number...: 4242424242424242
`
	require.Equal(t, expectedOutput, out.String())
}
