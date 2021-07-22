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

var Builder MessageSpecBuilder = &messageSpecBuilder{}

type MessageSpecBuilder interface {
	ImportJSON([]byte) (*iso8583.MessageSpec, error)
	ExportJSON(spec *iso8583.MessageSpec) ([]byte, error)
}

type messageSpecBuilder struct {
}

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
		_, err := fmt.Sscanf(key, "Field%d", &index)
		if err != nil {
			return nil, fmt.Errorf("invalid field index, index's format is `Field`+index ")
		}

		enc := getEncodingInterface(dummyField.Enc)
		pref := getPrefixInterface(dummyField.Prefix, dummyField.PrefixLength)
		pad := getPadInterface(dummyField.Padding)

		var newField field.Field

		switch dummyField.Type {
		case reflect.TypeOf(field.String{}).String():
			newField = field.NewString(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		case reflect.TypeOf(field.Numeric{}).String():
			newField = field.NewNumeric(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		case reflect.TypeOf(field.Binary{}).String():
			newField = field.NewBinary(&field.Spec{
				Length:      dummyField.Length,
				IDLength:    dummyField.IDLength,
				Description: dummyField.Description,
				Enc:         enc,
				Pref:        pref,
				Pad:         pad,
			})
		case reflect.TypeOf(field.Bitmap{}).String():
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
		dummyName := fmt.Sprintf("Field%03d", index)
		dummyField := fieldDummy{
			Type: reflect.TypeOf(orgField).Elem().String(),
		}

		switch reflect.TypeOf(orgField).Elem() {
		case reflect.TypeOf(field.String{}), reflect.TypeOf(field.Numeric{}), reflect.TypeOf(field.Binary{}), reflect.TypeOf(field.Bitmap{}):
			spec := orgField.Spec()

			dummyField.Length = spec.Length
			dummyField.IDLength = spec.IDLength
			dummyField.Description = spec.Description
			if spec.Enc != nil {
				dummyField.Enc = reflect.TypeOf(spec.Enc).Elem().String()
			}
			if spec.Pref != nil {
				dummyField.Prefix = reflect.TypeOf(spec.Pref).Elem().String()

				lengthStr := ""
				dummy := ""
				fmt.Sscanf(spec.Pref.Inspect(), "%s %s %s", &dummy, &lengthStr, &dummy)
				if lengthStr != "" && lengthStr != "fixed" && strings.ToUpper(lengthStr) == strings.Repeat("L", len(lengthStr)) {
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

	var pref prefix.Prefixer

	switch "*" + prefixName {
	case reflect.TypeOf(prefix.ASCII.Fixed).String():
		pref = prefix.ASCII.Fixed
	case reflect.TypeOf(prefix.ASCII.L).String():
		if length == 1 {
			pref = prefix.ASCII.L
		} else if length == 2 {
			pref = prefix.ASCII.LL
		} else if length == 3 {
			pref = prefix.ASCII.LLL
		} else if length == 4 {
			pref = prefix.ASCII.LLLL
		}
	case reflect.TypeOf(prefix.BCD.Fixed).String():
		pref = prefix.BCD.Fixed
	case reflect.TypeOf(prefix.BCD.L).String():
		if length == 1 {
			pref = prefix.BCD.L
		} else if length == 2 {
			pref = prefix.BCD.LL
		} else if length == 3 {
			pref = prefix.BCD.LLL
		} else if length == 4 {
			pref = prefix.BCD.LLLL
		}
	case reflect.TypeOf(prefix.Hex.Fixed).String():
		pref = prefix.Hex.Fixed
	case reflect.TypeOf(prefix.Hex.L).String():
		if length == 1 {
			pref = prefix.Hex.L
		} else if length == 2 {
			pref = prefix.Hex.LL
		} else if length == 3 {
			pref = prefix.Hex.LLL
		} else if length == 4 {
			pref = prefix.Hex.LLLL
		}
	case reflect.TypeOf(prefix.Binary.Fixed).String():
		pref = prefix.Binary.Fixed
	case reflect.TypeOf(prefix.Binary.L).String():
		if length == 1 {
			pref = prefix.Binary.L
		} else if length == 2 {
			pref = prefix.Binary.LL
		} else if length == 3 {
			pref = prefix.Binary.LLL
		} else if length == 4 {
			pref = prefix.Binary.LLLL
		}
	case reflect.TypeOf(prefix.EBCDIC.Fixed).String():
		pref = prefix.EBCDIC.Fixed
	case reflect.TypeOf(prefix.EBCDIC.L).String():
		if length == 1 {
			pref = prefix.EBCDIC.L
		} else if length == 2 {
			pref = prefix.EBCDIC.LL
		} else if length == 3 {
			pref = prefix.EBCDIC.LLL
		} else if length == 4 {
			pref = prefix.EBCDIC.LLLL
		}
	}

	return pref
}

func getEncodingInterface(encName string) encoding.Encoder {

	var enc encoding.Encoder

	switch "*" + encName {
	case reflect.TypeOf(encoding.ASCII).String():
		enc = encoding.ASCII
	case reflect.TypeOf(encoding.BCD).String():
		enc = encoding.BCD
	case reflect.TypeOf(encoding.Binary).String():
		enc = encoding.Binary
	case reflect.TypeOf(encoding.EBCDIC).String():
		enc = encoding.EBCDIC
	case reflect.TypeOf(encoding.Hex).String():
		enc = encoding.Hex
	case reflect.TypeOf(encoding.LBCD).String():
		enc = encoding.LBCD
	}
	return enc
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
