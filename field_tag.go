package iso8583

import (
	"reflect"
	"strconv"
	"strings"
)

type FieldTag struct {
	Id    int // is -1 if index is not a number
	Index string

	// KeepZero tells the marshaler to use zero value and set bitmap bit to
	// 1 for this field. Default behavior is to omit the field from the
	// message if it's zero value.
	KeepZero bool
}

func NewFieldTag(field reflect.StructField) FieldTag {
	// value of the key "index" in the tag
	var value string

	// keep the order of tags for now, when index tag is deprecated we can
	// change the order
	for _, tag := range []string{"index", "iso8583"} {
		if value = field.Tag.Get(tag); value != "" {
			break
		}
	}

	// format of the value is "id[,keep_zero_value]"
	// id is the id of the field
	// let's parse it
	if value != "" {
		index, keepZero := parseValue(value)

		id, err := strconv.Atoi(index)
		if err != nil {
			id = -1
		}

		return FieldTag{
			Id:       id,
			Index:    index,
			KeepZero: keepZero,
		}
	}

	dataFieldName := field.Name
	if len(dataFieldName) > 0 && fieldNameIndexRe.MatchString(dataFieldName) {
		indexStr := dataFieldName[1:]
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return FieldTag{
				Id:    -1,
				Index: indexStr,
			}
		}

		return FieldTag{
			Id:    fieldIndex,
			Index: indexStr,
		}
	}

	return FieldTag{
		Id: -1,
	}
}

func parseValue(value string) (index string, keepZero bool) {
	if value == "" {
		return
	}

	// split the value by comma
	values := strings.Split(value, ",")

	// the first value is the index
	index = values[0]

	// if there is only one value, return
	if len(values) == 1 {
		return
	}

	// if the second value is "keep_zero_value", set the flag
	if values[1] == "keepzero" {
		keepZero = true
	}

	return
}
