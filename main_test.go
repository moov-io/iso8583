package iso8583

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/spec"
	"github.com/stretchr/testify/require"
)

func TestISO8583(t *testing.T) {
	specTest := &spec.MessageSpec{
		Fields: map[int]spec.Packer{
			0: spec.MTI(4, "Message Type Indicator", encoding.ASCII),
			1: spec.Bitmap(16, "Bitmap", encoding.ASCII),
			2: spec.NewField(19, "Primary Account Number", encoding.ASCII, nil, nil, nil),
			3: spec.NewField(19, "Primary Account Number", encoding.ASCII, nil, nil, nil),
			5: spec.NewField(19, "Primary Account Number", encoding.ASCII, nil, nil, nil),
		},
	}

	message := NewMessage()
	message.Field(0, "0123")
	message.Field(2, "424242424242")
	message.Field(3, "424242424242")
	message.Field(5, "424242424242")

	message.Spec(specTest)
	raw, err := message.Pack()
	require.NoError(t, err)
	require.NotNil(t, raw)
	fmt.Println(raw)

	// message = NewMessage()
	// message.Spec(specTest)
	// message.Unpack(raw)

	// require.Equal(t, "0100", message.GetString(0))
	// require.Equal(t, "424242424242", message.GetString(2))
}
