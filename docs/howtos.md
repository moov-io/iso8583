# Howtos

## Howto create custom packer and unpacker for a field

### Problem

The default behavior of the field packer and unpacker may not meet your requirements. For instance, you might need the length prefix to represent the length of the encoded data, not the field value. This is often necessary when using BCD or HEX encoding, where the field value's length differs from the encoded field value's length.

**Example Requirement:**

- Maximum length of the field: 9
- Field value encoding: BCD
- Length prefix: L (1 byte) representing the length of the encoded data
- Field value: "123"

### Default Behavior

Let's explore the default behavior of a Numeric field:

```go
f := field.NewNumeric(&field.Spec{
    Length:      9, // The max length of the field is 9 digits
    Description: "Amount",
    Enc:         encoding.BCD,
    Pref:        prefix.Binary.L,
})

f.SetValue(123)

packed, err := f.Pack()
require.NoError(t, err)

require.Equal(t, []byte{0x03, 0x01, 0x23}, packed)
```

By default, the length prefix contains the field value's length, which is 3 digits, resulting in a length prefix of 0x03.

### Custom Packer and Unpacker

Let's create a custom packer and unpacker for the Numeric field to pack the field value as BCD and set the length prefix to the length of the encoded field value.

```go
f := field.NewNumeric(&field.Spec{
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

f.SetValue(123)

packed, err = f.Pack()
require.NoError(t, err)

require.Equal(t, []byte{0x02, 0x01, 0x23}, packed)
```

Since 123 encoded in BCD is 0x01, 0x23, the length prefix is 0x02, indicating the length of the packed data is 2 bytes, not the field value's length which is 3 digits.
