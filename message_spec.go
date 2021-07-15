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

func createMessageField(specField field.Field) field.Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to field.Field interface
	fl := reflect.New(fieldType).Interface().(field.Field)
	fl.SetSpec(specField.Spec())

	return fl
}
