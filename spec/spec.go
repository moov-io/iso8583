package spec

// import (
// 	"github.com/moov-io/iso8583/encoding"
// 	"github.com/moov-io/iso8583/fields"
// 	"github.com/moov-io/iso8583/prefixer"
// )

// type Prefixer interface{}
// type Padder interface{}
// type Validator interface{}

// type Packer interface {
// 	Pack(f fields.Field) ([]byte, error)
// 	Unpack(raw []byte) (fields.Field, error)
// }

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

// type MessageSpec struct {
// 	Fields map[int]Packer
// }
