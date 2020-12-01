package spec

import (
	"fmt"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefixer"
	"github.com/moov-io/iso8583/utils"
)

type bitmapField struct {
	Description string
	Enc         encoding.Encoder
	Pref        prefixer.Prefixer
}

func (fd *bitmapField) Pack(data []byte) ([]byte, error) {
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
func (fd *bitmapField) Unpack(data []byte) ([]byte, int, error) {
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
	return &bitmapField{
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}
