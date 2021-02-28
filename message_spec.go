package iso8583

import (
	"errors"
	"reflect"

	"github.com/moov-io/iso8583/field"
)

type MessageSpec struct {
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

// Get field's index by identifier
func (s *MessageSpec) GetFieldIndex(identifier string) (int, error) {

	for key, f := range s.Fields {
		if f.Spec().GetIdentifier() == identifier {
			return key, nil
		}
	}

	return 0, errors.New("don't find any specification by identifier")
}

func createMessageField(specField field.Field) field.Field {
	fieldType := reflect.TypeOf(specField).Elem()

	// create new field and convert it to field.Field interface
	fl := reflect.New(fieldType).Interface().(field.Field)
	fl.SetSpec(specField.Spec())

	return fl
}
