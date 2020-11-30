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
	Pack(data []byte) ([]byte, error)
	Unpack(data []byte) ([]byte, error)
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

func (fd *fieldDefinition) Unpack(data []byte) ([]byte, error) {
	dataLen, err := fd.Pref.DecodeLength(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	start := fd.Pref.Length()
	end := fd.Pref.Length() + dataLen
	raw, err := fd.Enc.Decode(data[start:end])
	if err != nil {
		return nil, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	return raw, nil
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

// Unpack of the Bitmap field returns data of varied length
// if there is only primary bitmap (bit 1 is not set) we return only 8 bytes
// if secondary bitmap presents (bit 1 is set) we return 16 bytes
func (fd *bitmapFieldDefinition) Unpack(data []byte) ([]byte, error) {
	dataLen, err := fd.Pref.DecodeLength(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	start := fd.Pref.Length()
	end := fd.Pref.Length() + dataLen
	raw, err := fd.Enc.Decode(data[start:end])
	if err != nil {
		return nil, fmt.Errorf("Failed to unpack '%s': %v", fd.Description, err)
	}

	// for the bitmap
	// if bit 1 is set then return 16 bytes
	// if bit 1 is not then return 8 byte
	bitmap := utils.NewBitmapFromData(raw)

	if bitmap.IsSet(1) {
		return raw[:16], nil
	}

	return raw[:8], nil
}

func Bitmap(desc string, enc encoding.Encoder, pref prefixer.Prefixer) Packer {
	return &bitmapFieldDefinition{
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}

// import (
// 	"github.com/moov-io/iso8583/encoding"
// 	"github.com/moov-io/iso8583/fields"
// 	"github.com/moov-io/iso8583/prefixer"
// )

// type Prefixer interface{}
// type Padder interface{}
// type Validator interface{}

// type FieldPacker struct {
// 	length      int
// 	description string
// 	enc         encoding.Encoder
// 	pref        Prefixer
// 	pad         Padder
// 	val         Validator
// }

// func (fp *FieldPacker) Pack(f fields.Field) ([]byte, error) {
// 	return f.Bytes(), nil
// }

// func (fp *FieldPacker) Unpack(raw []byte) (fields.Field, error) {
// 	decoded, err := fp.enc.Decode(raw)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return fields.NewField(-1, decoded), nil
// }

// func NewField(length int, description string, enc encoding.Encoder, pref prefixer.Prefixer, pad Padder, val Validator) Packer {
// 	return &FieldPacker{
// 		length:      length,
// 		description: description,
// 		enc:         enc,
// 		pref:        pref,
// 		pad:         pad,
// 		val:         val,
// 	}
// }

// func MTI(length int, description string, enc encoding.Encoder) Packer {
// 	return &FieldPacker{
// 		length: length,
// 		enc:    enc,
// 	}
// }

// func Bitmap(length int, description string, enc encoding.Encoder) Packer {
// 	return &FieldPacker{
// 		length: length,
// 		enc:    enc,
// 	}
// }

// type Definition struct {
// 	Fields map[int]Packer
// }
