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
		Length:      5, // now, this field indicates the max length of the encoded field value 9/2+1 = 5
		Description: "Amount",
		Enc:         encoding.BCD,
		Pref:        prefix.Binary.L,
		// we define a custom packer here, which will encode the length of the packed data
		Packer: field.PackerFunc(func(data []byte, spec *field.Spec) ([]byte, error) {
			if spec.Pad != nil {
				data = spec.Pad.Pad(data, spec.Length)
			}

			packed, err := spec.Enc.Encode(data)
			if err != nil {
				return nil, fmt.Errorf("failed to encode content: %w", err)
			}

			// here is where we encode the length of the packed data, not the length of the value
			packedLength, err := spec.Pref.EncodeLength(spec.Length, len(packed))
			if err != nil {
				return nil, fmt.Errorf("failed to encode length: %w", err)
			}

			return append(packedLength, packed...), nil
		}),
		// we define a custom unpacker here, which will decode the length of the packed data
		Unpacker: field.UnpackerFunc(func(data []byte, spec *field.Spec) ([]byte, int, error) {
			dataLen, prefBytes, err := spec.Pref.DecodeLength(spec.Length, data)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to decode length: %w", err)
			}

			// dataLen here is the length of the packed data, not the length of the value
			// as we use BCD decoding, we have to multiply it by 2, as each BCD byte
			// represents 2 digits. If the number of digits is even, it will be prepended
			// with a 0. As the type of the field is Numeric, leading 0 will be removed
			// so we will have exactly the number of digits we need.
			raw, read, err := spec.Enc.Decode(data[prefBytes:], dataLen*2)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to decode content: %w", err)
			}

			if spec.Pad != nil {
				raw = spec.Pad.Unpad(raw)
			}

			return raw, read + prefBytes, nil
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
