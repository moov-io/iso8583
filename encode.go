package iso8583

import (
	"errors"
	"fmt"
	"reflect"
)

func Marshal(message *Message, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("data is not a pointer or nil")
	}

	// get the struct from the pointer
	dataStruct := rv.Elem()

	if dataStruct.Kind() != reflect.Struct {
		return errors.New("data is not a struct")
	}

	// iterate over struct fields
	for i := 0; i < dataStruct.NumField(); i++ {
		fieldIndex, err := getFieldIndex(dataStruct.Type().Field(i))
		if err != nil {
			return fmt.Errorf("getting field %d index: %w", i, err)
		}

		// skip field without index
		if fieldIndex < 0 {
			continue
		}

		messageField := message.GetField(fieldIndex)
		// if struct field we are usgin to populate value expects to
		// set index of the field that is not described by spec
		if messageField == nil {
			return fmt.Errorf("no message field defined by spec with index: %d", fieldIndex)
		}

		dataField := dataStruct.Field(i)
		if dataField.IsNil() {
			continue
		}

		err = messageField.MarshalValue(dataField.Interface())
		if err != nil {
			return fmt.Errorf("failed to get data from field %d: %w", fieldIndex, err)
		}
	}

	return nil
}
