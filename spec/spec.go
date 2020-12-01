package spec

import (
	"fmt"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/moov-io/iso8583/utils"
)

type MessageSpec struct {
	Fields map[int]Packer
}

type Packer interface {
	// Pack packs data taking into account data encoding and data length
	// it returns packed data
	Pack(data []byte) ([]byte, error)

	// Unpack unpacks data taking into account data encoding and data length
	// it returns unpacked data and the number of bytes read
	Unpack(data []byte) ([]byte, int, error)
}

type fieldDefinition struct {
	Description string
	Enc         encoding.Encoder
	Pref        prefixer.Prefixer
}

func (fd *fieldDefinition) Pack(data []byte) ([]byte, error) {
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

func (fd *fieldDefinition) Unpack(data []byte) ([]byte, int, error) {
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
	return &fieldDefinition{
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}

type bitmapFieldDefinition struct {
	Description string
	Enc         encoding.Encoder
	Pref        prefixer.Prefixer
}

func (fd *bitmapFieldDefinition) Pack(data []byte) ([]byte, error) {
	bitmap := utils.NewBitmapFromData(data)

	packed, err := fd.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", fd.Description, err)
	}

	packedLength, err := fd.Pref.EncodeLength(len(packed))
	if err != nil {
		return nil, fmt.Errorf("Failed to pack '%s': %v", fd.Description, err)
	}

	if !bitmap.IsSet(1) {
		packed = packed[:len(packed)/2]
	}

	return append(packedLength, packed...), nil
}

// Unpack of the Bitmap field returns data of varied length
// if there is only primary bitmap (bit 1 is not set) we return only 8 bytes
// if secondary bitmap presents (bit 1 is set) we return 16 bytes
func (fd *bitmapFieldDefinition) Unpack(data []byte) ([]byte, int, error) {
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

	bitmap := utils.NewBitmapFromData(raw)

	if bitmap.IsSet(1) {
		return raw[:16], dataLen, nil
	}

	return raw[:8], dataLen / 2, nil
}

func Bitmap(desc string, enc encoding.Encoder, pref prefixer.Prefixer) Packer {
	return &bitmapFieldDefinition{
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}
