package iso8583

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// Unmarshal traverses the message fields recursively and for each field, sets
// the field value into the corresponding struct field value pointed by v. If v
// is nil or not a pointer it returns error.
func Unmarshal(message *Message, v interface{}) error {
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
		dataFieldName := dataStruct.Type().Field(i).Name

		// skip struct field if its name starts not from F
		if dataFieldName[0:1] != "F" {
			continue
		}

		indexStr := dataFieldName[1:]
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "converting field intex into int")
		}

		// we can get data only if field value is set
		messageField := message.GetField(fieldIndex)
		if messageField == nil {
			continue
		}

		dataField := dataStruct.Field(i)
		if dataField.IsNil() {
			dataField.Set(reflect.New(dataField.Type().Elem()))
		}

		err = messageField.UnmarshalValue(dataField.Interface())
		if err != nil {
			return fmt.Errorf("failed to get data from field %d: %w", fieldIndex, err)
		}
	}

	return nil
}
