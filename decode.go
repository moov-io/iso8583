package iso8583

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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
		fieldIndex, err := getFieldIndex(dataStruct.Type().Field(i))
		if err != nil {
			return fmt.Errorf("getting field %d index: %w", i, err)
		}

		// skip field without index
		if fieldIndex < 0 {
			continue
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

var indexFieldNameRe = regexp.MustCompile(`^F\d+$`)

// fieldIndex returns index of the field. First, it checks field name. If it
// does not match FNN (when NN is digits), it checks value of `index` tag.
// If negative value returned (-1) then index was not found for the field.
func getFieldIndex(field reflect.StructField) (int, error) {
	dataFieldName := field.Name

	if len(dataFieldName) > 0 && indexFieldNameRe.MatchString(dataFieldName) {
		indexStr := dataFieldName[1:]
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return -1, fmt.Errorf("converting field intex into int: %w", err)
		}

		return fieldIndex, nil
	}

	return -1, nil
}
