package field

import "fmt"

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
