package spec

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/stretchr/testify/require"
)

func TestSpec(t *testing.T) {
	spec87 := Definition{
		Fields: map[int]Packer{ // map[string]spec.Field
			// "0": spec.NewField(4, "MTI", encoder.ASCII, prefixer.None, padder.None, validator.Numeric),
			0: MTI(4, "Message Type Indicator", encoding.ASCII),
			1: Bitmap(16, "Bitmap", encoder.ASCII),
			// "2": spec.NewField(19, "Primary Account Number", encoder.ASCII, prefixer.RBCDLL, padder.No, validator.Numeric),
			// "3": spec.NewField(14, "Transaction amount", encoder.ASCII, prefixer.None, padder.LeftZero, validator.Numeric),

			// // or pre-built fields
			// "3": spec.ALLNumeric(19, "Primary Account Number"),
			// "4": spec.ANumeric(14, "Transaction amount"),
			// ...
		},
	}

	fmt.Println(spec87)

	// packer := MTI(4, "MTI", ASCII)
	// field, err := packer.Unpack([]byte("0100"))
}

func TestFieldPacker(t *testing.T) {
	t.Run("Unpack ASCII", func(t *testing.T) {
		packer := NewField(4, "MTI", encoding.ASCII, nil, nil, nil)

		field, err := packer.Unpack([]byte("0100"))

		require.NoError(t, err)
		require.NotNil(t, field)
		require.Equal(t, "0100", field.String())
	})

	t.Run("Unpack Binary", func(t *testing.T) {
		packer := NewField(4, "MTI", encoding.BCD, nil, nil, nil)

		field, err := packer.Unpack([]byte{0x01, 0x00})

		require.NoError(t, err)
		require.NotNil(t, field)
		require.Equal(t, "0100", field.String())
	})
}
