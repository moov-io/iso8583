package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

const (
	FieldPrefix = "Field"
)

type fieldConstructorFunc func(spec *field.Spec) field.Field

var (
	fieldConstructor = map[string]fieldConstructorFunc{
		"String":  newString,
		"Numeric": newNumeric,
		"Binary":  newBinary,
		"Bitmap":  newBitmap,
	}
	prefixes = map[string]prefix.Prefixer{
		"ascii.fixed":  prefix.ASCII.Fixed,
		"ascii.l":      prefix.ASCII.L,
		"ascii.ll":     prefix.ASCII.LL,
		"ascii.lll":    prefix.ASCII.LLL,
		"ascii.llll":   prefix.ASCII.LLLL,
		"bcd.fixed":    prefix.BCD.Fixed,
		"bcd.l":        prefix.BCD.L,
		"bcd.ll":       prefix.BCD.LL,
		"bcd.lll":      prefix.BCD.LLL,
		"bcd.llll":     prefix.BCD.LLLL,
		"hex.fixed":    prefix.Hex.Fixed,
		"hex.l":        prefix.Hex.L,
		"hex.ll":       prefix.Hex.LL,
		"hex.lll":      prefix.Hex.LLL,
		"hex.llll":     prefix.Hex.LLLL,
		"ebcdic.fixed": prefix.EBCDIC.Fixed,
		"ebcdic.l":     prefix.EBCDIC.L,
		"ebcdic.ll":    prefix.EBCDIC.LL,
		"ebcdic.lll":   prefix.EBCDIC.LLL,
		"ebcdic.llll":  prefix.EBCDIC.LLLL,
		"binary.fixed": prefix.Binary.Fixed,
	}
	encodings = map[string]encoding.Encoder{
		"asciiEncoder":  encoding.ASCII,
		"bcdEncoder":    encoding.BCD,
		"ebcdicEncoder": encoding.EBCDIC,
		"binaryEncoder": encoding.Binary,
		"hexEncoder":    encoding.Hex,
		"lBCDEncoder":   encoding.LBCD,
	}
)

var Builder MessageSpecBuilder = &messageSpecBuilder{}

type MessageSpecBuilder interface {
	ImportJSON([]byte) (*iso8583.MessageSpec, error)
	ExportJSON(spec *iso8583.MessageSpec) ([]byte, error)
}

type messageSpecBuilder struct{}

type specDummy struct {
	Name   string                `json:"name,omitempty" xml:"name,omitempty"`
	Fields map[string]fieldDummy `json:"fields,omitempty" xml:"fields,omitempty"`
}

type fieldDummy struct {
	Type        string    `json:"type,omitempty" xml:"type,omitempty"`
	Length      int       `json:"length,omitempty" xml:"length,omitempty"`
	IDLength    int       `json:"id_length,omitempty" xml:"id_length,omitempty"`
	Description string    `json:"description,omitempty" xml:"description,omitempty"`
	Enc         string    `json:"enc,omitempty" xml:"enc,omitempty"`
	Prefix      string    `json:"prefix,omitempty" xml:"prefix,omitempty"`
	Padding     *padDummy `json:"padding,omitempty" xml:"padding,omitempty"`
}

type padDummy struct {
	Type string `json:"type" xml:"type"`
	Pad  rune   `json:"pad" xml:"pad"`
}

func (builder *messageSpecBuilder) ImportJSON(raw []byte) (*iso8583.MessageSpec, error) {
	dummy := specDummy{}
	err := json.Unmarshal(raw, &dummy)
	if err != nil {
		return nil, err
	}

	if len(dummy.Fields) == 0 {
		return nil, fmt.Errorf("invalid json spec file")
	}

	spec := iso8583.MessageSpec{
		Name:   dummy.Name,
		Fields: make(map[int]field.Field),
	}

	for key, dummyField := range dummy.Fields {
		index := 0
		_, err := fmt.Sscanf(key, FieldPrefix+"%d", &index)
		if err != nil {
			return nil, fmt.Errorf("invalid field index, index's format is `Field`+index ")
		}

		constructor := fieldConstructor[dummyField.Type]
		if constructor == nil {
			return nil, fmt.Errorf("unable create field with %s", dummyField.Type)
		}

		enc := encodings[dummyField.Enc]
		if enc == nil {
			return nil, fmt.Errorf("encoding(%s) is incorrect, enc is mandatory field", dummyField.Enc)
		}
		pref := prefixes[strings.ToLower(dummyField.Prefix)]
		if pref == nil {
			return nil, fmt.Errorf("prefix(%s) is incorrect, prefix is mandatory field", dummyField.Prefix)
		}

		pad := getPadInterface(dummyField.Padding)

		spec.Fields[index] = constructor(&field.Spec{
			Length:      dummyField.Length,
			IDLength:    dummyField.IDLength,
			Description: dummyField.Description,
			Enc:         enc,
			Pref:        pref,
			Pad:         pad,
		})

	}

	return &spec, nil
}

func (builder *messageSpecBuilder) ExportJSON(orgSpec *iso8583.MessageSpec) ([]byte, error) {

	if orgSpec == nil {
		return nil, fmt.Errorf("invalid message spec")
	}

	dummyJson := specDummy{
		Name:   orgSpec.Name,
		Fields: make(map[string]fieldDummy),
	}

	for index, orgField := range orgSpec.Fields {
		dummyName := fmt.Sprintf(FieldPrefix+"%03d", index)
		dummyField := fieldDummy{
			Type: reflect.TypeOf(orgField).Elem().Name(),
		}

		switch reflect.TypeOf(orgField).Elem() {
		case reflect.TypeOf(field.String{}), reflect.TypeOf(field.Numeric{}), reflect.TypeOf(field.Binary{}), reflect.TypeOf(field.Bitmap{}):
			spec := orgField.Spec()

			dummyField.Length = spec.Length
			dummyField.IDLength = spec.IDLength
			dummyField.Description = spec.Description
			if spec.Enc != nil {
				dummyField.Enc = reflect.TypeOf(spec.Enc).Elem().Name()
			}
			if spec.Pref != nil {
				var lengthStr, prefixStr, str2 string
				if fmt.Sscanf(spec.Pref.Inspect(), "%s %s %s", &prefixStr, &lengthStr, &str2); len(prefixStr) > 0 {
					dummyField.Prefix = prefixStr + "." + lengthStr
				}
			}
			if spec.Pad != nil {
				pad := padDummy{}
				switch reflect.TypeOf(spec.Pad).Elem().String() {
				case "padding.leftPadder":
					if spec.Pad.Inspect() != nil {
						pad.Type = reflect.TypeOf(spec.Pad).Elem().Name()
						pad.Pad, _ = utf8.DecodeRune(spec.Pad.Inspect())
					}
				}
				if pad.Type != "" {
					dummyField.Padding = &pad
				}
			}

		default:
			continue
		}

		dummyJson.Fields[dummyName] = dummyField
	}

	outputBuffer := new(bytes.Buffer)
	enc := json.NewEncoder(outputBuffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err := enc.Encode(dummyJson)

	if err != nil {
		return nil, fmt.Errorf("unable to export message spec")
	}

	return outputBuffer.Bytes(), nil
}

func getPadInterface(info *padDummy) padding.Padder {
	if info == nil || info.Type == "" {
		return nil
	}

	var pad padding.Padder
	switch info.Type {
	case "leftPadder":
		pad = padding.Left(info.Pad)
	}

	return pad
}

func newString(spec *field.Spec) field.Field { return field.NewString(spec) }

func newNumeric(spec *field.Spec) field.Field { return field.NewNumeric(spec) }

func newBinary(spec *field.Spec) field.Field { return field.NewBinary(spec) }

func newBitmap(spec *field.Spec) field.Field { return field.NewBitmap(spec) }
