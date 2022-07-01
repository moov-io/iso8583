package iso8583

import (
	"reflect"

	"github.com/moov-io/iso8583/field"
)

type MessageSpec struct {
	Name   string
	Fields map[int]field.Field
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

// GetBitmapsSize returns size of Bit Maps based on spec
// Min : 1
// Max : 10
func (s *MessageSpec) GetBitmapsSize() (size int) {

	size = 1 // primary
	maxIndex := 0

	for key, elm := range s.Fields {
		if elm.Spec() != nil && key > maxIndex {
			maxIndex = key
		}
	}

	if maxIndex < 65 {
		return
	}

	remainder := maxIndex % 64
	if remainder != 0 {
		remainder = 1
	}

	size = maxIndex/64 + remainder

	if size > 10 { // max check
		size = 10
	}

	return
}

func createMessageField(specField field.Field) field.Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to field.Field interface
	fl := reflect.New(fieldType).Interface().(field.Field)
	fl.SetSpec(specField.Spec())

	// if it's a composite field, we have to recusively create its subfields as well
	if composite, ok := fl.(field.CompositeWithSubfields); ok {
		composite.ConstructSubfields()
	}

	return fl
}
