package field

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var fieldNameIndexRe = regexp.MustCompile(`^F.+$`)

type IndexTag struct {
	Id int // is -1 if index is not a number

	Tag string
	// KeepZero tells the marshaler to use zero value and set bitmap bit to
	// 1 for this field. Default behavior is to omit the field from the
	// message if it's zero value.
	KeepZero bool
}

func NewIndexTag(field reflect.StructField) IndexTag {
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
		tag, opts := parseTag(value)

		id, err := strconv.Atoi(tag)
		if err != nil {
			id = -1
		}

		return IndexTag{
			Id:       id,
			Tag:      tag,
			KeepZero: opts.Contains("keepzero"),
		}
	}

	dataFieldName := field.Name
	if len(dataFieldName) > 0 && fieldNameIndexRe.MatchString(dataFieldName) {
		indexStr := dataFieldName[1:]
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return IndexTag{
				Id:  -1,
				Tag: indexStr,
			}
		}

		return IndexTag{
			Id:  fieldIndex,
			Tag: indexStr,
		}
	}

	return IndexTag{
		Id: -1,
	}
}

type tagOptions string

// parseTag splits a struct field's index tag into its id and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}
