package field

import (
	"github.com/stoewer/go-strcase"

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
	Identifier  string
}

func (s *Spec) GetIdentifier() string {
	var identifier string
	if len(s.Identifier) == 0 {
		identifier = strcase.SnakeCase(s.Description)
	} else {
		identifier = strcase.SnakeCase(s.Identifier)
	}
	return identifier
}

func NewSpec(id string, length int, desc string, enc encoding.Encoder, pref prefix.Prefixer) *Spec {
	return &Spec{
		Identifier:  id,
		Length:      length,
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}
