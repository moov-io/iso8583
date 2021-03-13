package field

import (
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

type Spec struct {
	Length      int
	Description string
	Enc         encoding.Encoder
	Pref        prefix.Prefixer
	Pad         padding.Padder
}

func NewSpec(length int, desc string, enc encoding.Encoder, pref prefix.Prefixer) *Spec {
	return &Spec{
		Length:      length,
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}
