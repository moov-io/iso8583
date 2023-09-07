package iso8583

import (
	"fmt"
	"reflect"

	"github.com/moov-io/iso8583/field"
)

type MessageSpec struct {
	Name   string
	Fields map[int]field.Field
}

// Validate checks if the MessageSpec is valid.
func (s *MessageSpec) Validate() error {
	// we require MTI and Bitmap fields
	if _, ok := s.Fields[mtiIdx]; !ok {
		return fmt.Errorf("MTI field (%d) is required", mtiIdx)
	}

	if _, ok := s.Fields[bitmapIdx]; !ok {
		return fmt.Errorf("Bitmap field (%d) is required", bitmapIdx)
	}

	// check type of the bitmap field
	if _, ok := s.Fields[bitmapIdx].(*field.Bitmap); !ok {
		return fmt.Errorf("Bitmap field (%d) must be of type *field.Bitmap", bitmapIdx)
	}

	return nil
}

// Creates a map with new instances of Fields (Field interface)
// based on the field type in MessageSpec.
func (s *MessageSpec) CreateMessageFields() map[int]field.Field {

	fields := map[int]field.Field{}

	for k, specField := range s.Fields {
		fields[k] = createMessageField(specField)
	}

	return fields
}

func createMessageField(specField field.Field) field.Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to field.Field interface

	//nolint:forcetypeassert // we know the type of the field we're creating here
	fl := reflect.New(fieldType).Interface().(field.Field)
	fl.SetSpec(specField.Spec())

	// if it's a composite field, we have to recusively create its subfields as well
	if composite, ok := fl.(field.CompositeWithSubfields); ok {
		composite.ConstructSubfields()
	}

	return fl
}
