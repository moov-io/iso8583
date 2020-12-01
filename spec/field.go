package spec

import (
	"fmt"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
)

type field struct {
	Description string
	Enc         encoding.Encoder
	Pref        prefixer.Prefixer
}

func (fd *field) Pack(data []byte) ([]byte, error) {
	packed, err := fd.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", fd.Description, err)
	}

	packedLength, err := fd.Pref.EncodeLength(len(packed))
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", fd.Description, err)
	}

	return append(packedLength, packed...), nil
}

func (fd *field) Unpack(data []byte) ([]byte, int, error) {
	dataLen, err := fd.Pref.DecodeLength(data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	start := fd.Pref.Length()
	end := fd.Pref.Length() + dataLen
	raw, err := fd.Enc.Decode(data[start:end])
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	return raw, dataLen + fd.Pref.Length(), nil
}

func NewField(desc string, enc encoding.Encoder, pref prefixer.Prefixer) Packer {
	return &field{
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}
