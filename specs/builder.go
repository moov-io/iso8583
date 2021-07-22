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
	Type         string    `json:"type,omitempty" xml:"type,omitempty"`
	Length       int       `json:"length,omitempty" xml:"length,omitempty"`
	IDLength     int       `json:"id_length,omitempty" xml:"id_length,omitempty"`
	Description  string    `json:"description,omitempty" xml:"description,omitempty"`
	Enc          string    `json:"enc,omitempty" xml:"enc,omitempty"`
	Prefix       string    `json:"prefix,omitempty" xml:"prefix,omitempty"`
	PrefixLength int       `json:"prefix_length,omitempty" xml:"prefix_length,omitempty"`
	Padding      *padDummy `json:"padding,omitempty" xml:"padding,omitempty"`
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

		enc := getEncodingInterface(dummyField.Enc)
		pref := getPrefixInterface(dummyField.Prefix, dummyField.PrefixLength)
		pad := getPadInterface(dummyField.Padding)

		var newField field.Field

		switch dummyField.Type {
		case reflect.TypeOf(field.String{}).Name():
			newField = field.NewString(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		case reflect.TypeOf(field.Numeric{}).Name():
			newField = field.NewNumeric(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		case reflect.TypeOf(field.Binary{}).Name():
			newField = field.NewBinary(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		case reflect.TypeOf(field.Bitmap{}).Name():
			newField = field.NewBitmap(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		default:
			continue
		}

		spec.Fields[index] = newField
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
				dummyField.Prefix = reflect.TypeOf(spec.Pref).Elem().Name()

				var lengthStr, str1, str2 string
				if fmt.Sscanf(spec.Pref.Inspect(), "%s %s %s", &str1, &lengthStr, &str2); len(lengthStr) > 0 && lengthStr != "fixed" && strings.ToUpper(lengthStr) == strings.Repeat("L", len(lengthStr)) {
					dummyField.PrefixLength = len(lengthStr)
				}
			}
			if spec.Pad != nil {
				pad := padDummy{}
				switch reflect.TypeOf(spec.Pad).Elem().String() {
				case "padding.leftPadder":
					if spec.Pad.Inspect() != nil {
						pad.Type = reflect.TypeOf(spec.Pad).Elem().String()
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

func getPrefixInterface(prefixName string, length int) prefix.Prefixer {

	var prefixes = map[string]prefix.Prefixer{
		"asciiFixedPrefixer":     prefix.ASCII.Fixed,
		"asciiVarPrefixer.L":     prefix.ASCII.L,
		"asciiVarPrefixer.LL":    prefix.ASCII.LL,
		"asciiVarPrefixer.LLL":   prefix.ASCII.LLL,
		"asciiVarPrefixer.LLLL":  prefix.ASCII.LLLL,
		"bcdFixedPrefixer":       prefix.BCD.Fixed,
		"bcdVarPrefixer.L":       prefix.BCD.L,
		"bcdVarPrefixer.LL":      prefix.BCD.LL,
		"bcdVarPrefixer.LLL":     prefix.BCD.LLL,
		"bcdVarPrefixer.LLLL":    prefix.BCD.LLLL,
		"hexFixedPrefixer":       prefix.Hex.Fixed,
		"hexVarPrefixer.L":       prefix.Hex.L,
		"hexVarPrefixer.LL":      prefix.Hex.LL,
		"hexVarPrefixer.LLL":     prefix.Hex.LLL,
		"hexVarPrefixer.LLLL":    prefix.Hex.LLLL,
		"ebcdicFixedPrefixer":    prefix.EBCDIC.Fixed,
		"ebcdicVarPrefixer.L":    prefix.EBCDIC.L,
		"ebcdicVarPrefixer.LL":   prefix.EBCDIC.LL,
		"ebcdicVarPrefixer.LLL":  prefix.EBCDIC.LLL,
		"ebcdicVarPrefixer.LLLL": prefix.EBCDIC.LLLL,
		"binaryFixedPrefixer":    prefix.Binary.Fixed,
	}

	if length > 0 {
		prefixName = prefixName + "." + strings.Repeat("L", length)
	}

	return prefixes[prefixName]
}

func getEncodingInterface(encName string) encoding.Encoder {

	var encodes = map[string]encoding.Encoder{
		"asciiEncoder":  encoding.ASCII,
		"bcdEncoder":    encoding.BCD,
		"ebcdicEncoder": encoding.EBCDIC,
		"binaryEncoder": encoding.Binary,
		"hexEncoder":    encoding.Hex,
		"lBCDEncoder":   encoding.LBCD,
	}

	return encodes[encName]
}

func getPadInterface(info *padDummy) padding.Padder {
	if info == nil || info.Type == "" {
		return nil
	}

	var pad padding.Padder
	switch info.Type {
	case "padding.leftPadder":
		pad = padding.Left(info.Pad)
	}

	return pad
}
