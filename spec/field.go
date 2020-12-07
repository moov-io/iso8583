package spec

import (
	"fmt"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

type Field struct {
	Length      int
	Description string
	Enc         encoding.Encoder
	Pref        prefix.Prefixer
	Pad         padding.Padder
}

func NewField(length int, desc string, enc encoding.Encoder, pref prefix.Prefixer) Packer {
	return &Field{
		Length:      length,
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}

func (fd *Field) Pack(data []byte) ([]byte, error) {
	if fd.Pad != nil {
		data = fd.Pad.Pad(data, fd.Length)
	}

	packed, err := fd.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", fd.Description, err)
	}

	packedLength, err := fd.Pref.EncodeLength(fd.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", fd.Description, err)
	}

	return append(packedLength, packed...), nil
}

func (fd *Field) Unpack(data []byte) ([]byte, int, error) {
	dataLen, err := fd.Pref.DecodeLength(fd.Length, data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	start := fd.Pref.Length()
	end := fd.Pref.Length() + dataLen
	raw, err := fd.Enc.Decode(data[start:end])
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	if fd.Pad != nil {
		raw = fd.Pad.Unpad(raw)
	}

	return raw, dataLen + fd.Pref.Length(), nil
}
