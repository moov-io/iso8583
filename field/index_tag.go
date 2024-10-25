package field

import (
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var fieldNameIndexRe = regexp.MustCompile(`^F.+$`)

type IndexTag struct {
	ID int // is -1 if index is not a number

	Tag string
	// KeepZero tells the marshaler to use zero value and set bitmap bit to
	// 1 for this field. Default behavior is to omit the field from the
	// message if it's zero value.
	KeepZero bool
}

func NewIndexTag(field reflect.StructField) IndexTag {
	indexTag := extractTagInfo(field)

	if indexTag.Tag == "" {
		id, tag := extractIdAndTagFromName(field.Name)
		indexTag.ID = id
		indexTag.Tag = tag
	}

	return indexTag
}

func extractTagInfo(field reflect.StructField) IndexTag {
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
	if value == "" {
		return IndexTag{
			ID: -1,
		}
	}

	tag, opts := parseTag(value)
	id, err := strconv.Atoi(tag)
	if err != nil {
		id = -1
	}

	return IndexTag{
		ID:       id,
		Tag:      tag,
		KeepZero: opts.Contains("keepzero"),
	}
}

func extractIdAndTagFromName(fieldName string) (int, string) {
	if len(fieldName) > 0 && fieldNameIndexRe.MatchString(fieldName) {
		indexStr := fieldName[1:]
		fieldIndex, err := strconv.Atoi(indexStr)
		if err != nil {
			return -1, indexStr
		}

		return fieldIndex, indexStr
	}

	return -1, ""
}

type tagOptions []string

// parseTag splits a struct field's index tag into its id and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")

	return tag, tagOptions(strings.Split(opt, ","))
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	return slices.Contains(o, optionName)
}
