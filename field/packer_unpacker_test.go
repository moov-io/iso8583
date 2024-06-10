package field_test

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestCustomPackerAndUnpacker(t *testing.T) {
	// let's say we have a requirements that the lentgh prefix should
	// contain the length of the data we pass, not the length of the field
	// value. This is the case when you use BCD or HEX encoding, where the
	// length of the field value is not the same as the length of the
	// encoded field value.

	// Here is an example of such requirement:
	// - the max length of the field is 9
	// - the field value should be encoded using BCD encoding
	// - the lenth prefix is L (1 byte) and should contain the length of the
	//   the data in the field we pass.
	// - the field value is "123"

	// let's see the default behavior of the Numeric field
	fd := field.NewNumeric(&field.Spec{
		Length:      9, // the max length of the field is 9 digits
		Description: "Amount",
		Enc:         encoding.BCD,
		Pref:        prefix.Binary.L,
	})

	fd.SetValue(123)

	packed, err := fd.Pack()
	require.NoError(t, err)

	// we expect the length to be 2 bytes, as 123 encoded in BCD is 0x01, 0x23
	// by the default behavior, the length prefix will contain the length of the
	// field value, which is 3 digits, so the length prefix as you can see is 0x03
	require.Equal(t, []byte{0x03, 0x01, 0x23}, packed)

	// now let's create a custom packer and unpacker for the Numeric field
	// that will pack the field value as BCD and the length prefix as the length
	// of the encoded field value.
	fc := field.NewNumeric(&field.Spec{
		Length:      9, // max length of the field value (9 digits)
		Description: "Amount",
		Enc:         encoding.BCD,
		Pref:        prefix.Binary.L,
		// Define a custom packer to encode the length of the packed data
		Packer: field.PackerFunc(func(value []byte, spec *field.Spec) ([]byte, error) {
			if spec.Pad != nil {
				value = spec.Pad.Pad(value, spec.Length)
			}

			encodedValue, err := spec.Enc.Encode(value)
			if err != nil {
				return nil, fmt.Errorf("failed to encode content: %w", err)
			}

			// Encode the length of the packed data, not the length of the value
			maxLength := spec.Length/2 + 1

			// Encode the length of the encoded value
			lengthPrefix, err := spec.Pref.EncodeLength(maxLength, len(encodedValue))
			if err != nil {
				return nil, fmt.Errorf("failed to encode length: %w", err)
			}

			return append(lengthPrefix, encodedValue...), nil
		}),

		// Define a custom unpacker to decode the length of the packed data
		Unpacker: field.UnpackerFunc(func(packedFieldValue []byte, spec *field.Spec) ([]byte, int, error) {
			maxEncodedValueLength := spec.Length/2 + 1

			encodedValueLength, prefBytes, err := spec.Pref.DecodeLength(maxEncodedValueLength, packedFieldValue)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to decode length: %w", err)
			}

			// for BCD encoding, the length of the packed data is twice the length of the encoded value
			valueLength := encodedValueLength * 2

			// Decode the packed data length
			value, read, err := spec.Enc.Decode(packedFieldValue[prefBytes:], valueLength)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to decode content: %w", err)
			}

			if spec.Pad != nil {
				value = spec.Pad.Unpad(value)
			}

			return value, read + prefBytes, nil
		}),
	})

	fc.SetValue(123)

	packed, err = fc.Pack()
	require.NoError(t, err)

	// we expect the length to be 2 bytes, as 123 encoded in BCD is 0x01, 0x23
	// so, you can see that the length prefix is 0x02, as the length of the packed
	// data is 2 bytes.
	require.Equal(t, []byte{0x02, 0x01, 0x23}, packed)
}
