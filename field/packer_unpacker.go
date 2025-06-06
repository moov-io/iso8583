package field

import (
	"fmt"
)

type defaultPacker struct{}

// Pack packs the data according to the spec
func (p defaultPacker) Pack(value []byte, spec *Spec) ([]byte, error) {
	// pad the value if needed
	if spec.Pad != nil {
		value = spec.Pad.Pad(value, spec.Length)
	}

	// encode the value
	encodedValue, err := spec.Enc.Encode(value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	// encode the length
	lengthPrefix, err := spec.Pref.EncodeLength(spec.Length, len(value))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(lengthPrefix, encodedValue...), nil
}

type defaultUnpacker struct{}

// Unpack unpacks the data according to the spec
func (u defaultUnpacker) Unpack(packedFieldValue []byte, spec *Spec) ([]byte, int, error) {
	// decode the length
	valueLength, prefBytes, err := spec.Pref.DecodeLength(spec.Length, packedFieldValue)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode length: %w", err)
	}

	// decode the value
	value, read, err := spec.Enc.Decode(packedFieldValue[prefBytes:], valueLength)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode content: %w", err)
	}

	// unpad the value if needed
	if spec.Pad != nil {
		value = spec.Pad.Unpad(value)
	}

	return value, read + prefBytes, nil
}

type Track2Packer struct{}

// This is a custom packer for Track2 Data. Some specifications require the length
// to be equal to the length of the pre-padded value.
func (p Track2Packer) Pack(value []byte, spec *Spec) ([]byte, error) {
	data := value

	// Only pad if the length is odd. If so, just add
	// one pad character, so tell the Pad function that
	// the length we want is +1 to what the value is
	if spec.Pad != nil && len(value)%2 != 0 {
		data = spec.Pad.Pad(data, len(value)+1)
	}

	packed, err := spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	// Encode the length to that of the original string, not the potentially
	// padded length
	packedLength, err := spec.Pref.EncodeLength(spec.Length, len(value))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

type Track2Unpacker struct{}

func (p Track2Unpacker) Unpack(packedFieldValue []byte, spec *Spec) ([]byte, int, error) {
	// decode the length
	valueLength, prefBytes, err := spec.Pref.DecodeLength(spec.Length, packedFieldValue)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode length: %w", err)
	}

	// if valueLength is odd we need to make it even to adjust for
	// the padding in our Packer
	if valueLength%2 != 0 {
		valueLength++
	}

	// decode the value
	value, read, err := spec.Enc.Decode(packedFieldValue[prefBytes:], valueLength/2)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode content: %w", err)
	}

	// unpad the value if needed
	if spec.Pad != nil {
		value = spec.Pad.Unpad(value)
	}

	return value, read + prefBytes, nil
}
