package spec

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/stretchr/testify/require"
)

// Encoders:
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

func TestSpec(t *testing.T) {
	// spec87 := Definition{
	// 	Fields: map[int]Packer{ // map[string]spec.Field
	// 		// "0": spec.NewField(4, "MTI", encoder.ASCII, prefixer.None, padder.None, validator.Numeric),
	// 		0: MTI(4, "Message Type Indicator", encoding.ASCII),
	// 		1: Bitmap(16, "Bitmap", encoder.ASCII),
	// 		// "2": spec.NewField(19, "Primary Account Number", encoder.ASCII, prefixer.RBCDLL, padder.No, validator.Numeric),
	// 		// "3": spec.NewField(14, "Transaction amount", encoder.ASCII, prefixer.None, padder.LeftZero, validator.Numeric),

	// 		// // or pre-built fields
	// 		// "3": spec.ALLNumeric(19, "Primary Account Number"),
	// 		// "4": spec.ANumeric(14, "Transaction amount"),
	// 		// ...
	// 	},
	// }

	// fmt.Println(spec87)

	// packer := MTI(4, "MTI", ASCII)
	// field, err := packer.Unpack([]byte("0100"))
}

type fieldDefinition struct {
	Description string
	Enc         encoding.Encoder
	Pref        prefixer.Prefixer
}

func (fd *fieldDefinition) Pack(data []byte) ([]byte, error) {

	packed, err := fd.Enc.Encode(data)
	if err != nil {
		return nil, err
	}

	packedLength, err := fd.Pref.EncodeLength(len(packed))
	if err != nil {
		return nil, err
	}

	return append(packedLength, packed...), nil
}

func (fd *fieldDefinition) Unpack(data []byte) ([]byte, error) {
	dataLen, err := fd.Pref.DecodeLength(data)
	if err != nil {
		return nil, err
	}

	start := fd.Pref.Length()
	end := fd.Pref.Length() + dataLen
	raw, err := fd.Enc.Decode(data[start:end])
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func TestFieldPacker(t *testing.T) {
	t.Run("ASCII VAR field", func(t *testing.T) {
		field := &fieldDefinition{"Primary Account Number", encoding.ASCII, prefixer.ASCII.LL(19)}

		// pack
		got, err := field.Pack([]byte("4242424242424242"))
		want := []byte{49, 54, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50}
		require.NoError(t, err)
		require.Equal(t, want, got)

		// unpack
		got, err = field.Unpack([]byte{49, 54, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50, 52, 50})
		want = []byte("4242424242424242")
		require.NoError(t, err)
		require.Equal(t, want, got)

	})

	t.Run("ASCII Fixed field", func(t *testing.T) {
		field := &fieldDefinition{"Processing code", encoding.ASCII, prefixer.ASCII.Fixed(6)}

		// pack
		got, err := field.Pack([]byte("123456"))
		want := []byte("123456")
		require.NoError(t, err)
		require.Equal(t, want, got)

		// unpack
		got, err = field.Unpack([]byte("123456"))
		want = []byte("123456")
		require.NoError(t, err)
		require.Equal(t, want, got)

	})
}
