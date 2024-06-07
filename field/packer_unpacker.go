package field

import "fmt"

type DefaultPacker struct{}

func (p DefaultPacker) Pack(data []byte, spec *Spec) ([]byte, error) {
	if spec.Pad != nil {
		data = spec.Pad.Pad(data, spec.Length)
	}

	packed, err := spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	packedLength, err := spec.Pref.EncodeLength(spec.Length, len(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

type DefaultUnpacker struct{}

func (u DefaultUnpacker) Unpack(data []byte, spec *Spec) ([]byte, int, error) {
	dataLen, prefBytes, err := spec.Pref.DecodeLength(spec.Length, data)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode length: %w", err)
	}

	raw, read, err := spec.Enc.Decode(data[prefBytes:], dataLen)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode content: %w", err)
	}

	if spec.Pad != nil {
		raw = spec.Pad.Unpad(raw)
	}

	return raw, read + prefBytes, nil
}
