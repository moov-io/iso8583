package iso8583

import (
	"fmt"
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

	// fmt.Println(fields)
	// fmt.Println(fields[0])
	// fmt.Println(reflect.TypeOf(fields[0]))

	// fmt.Println("type of the field:", reflect.TypeOf(fields[0]).Elem())

	return fields
}

func createMessageField(specField field.Field) field.Field {
	fmt.Println("hello")
	fieldType := reflect.TypeOf(specField).Elem()

	fmt.Println("type of the field:", fieldType)

	// create new field and convert it to field.Field interface
	fl := reflect.New(fieldType).Interface().(field.Field)
	fl.SetSpec(specField.Spec())

	fmt.Println("new field", fl)

	return fl
}
