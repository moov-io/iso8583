package field

import (
	"reflect"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

type Spec struct {
	Length int
	// Only applicable Composite and other bespoke field types.
	IDLength    int
	Description string
	Enc         encoding.Encoder
	Pref        prefix.Prefixer
	Pad         padding.Padder
	Fields      map[int]Field
	CountT      string
}

func NewSpec(length int, desc string, enc encoding.Encoder, pref prefix.Prefixer) *Spec {
	return &Spec{
		Length:      length,
		Description: desc,
		Enc:         enc,
		Pref:        pref,
	}
}

// Creates a map with new instances of Fields (Field interface)
// based on the field type in Spec.
func (s *Spec) CreateMessageFields() map[int]Field {
	fields := map[int]Field{}

	for k, specField := range s.Fields {
		fields[k] = createMessageField(specField)
	}

	return fields
}

func createMessageField(specField Field) Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to Field interface
	fl := reflect.New(fieldType).Interface().(Field)
	fl.SetSpec(specField.Spec())

	return fl
}
