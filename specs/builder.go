package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

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
		"String":  func(spec *field.Spec) field.Field { return field.NewString(spec) },
		"Numeric": func(spec *field.Spec) field.Field { return field.NewNumeric(spec) },
		"Binary":  func(spec *field.Spec) field.Field { return field.NewBinary(spec) },
		"Bitmap":  func(spec *field.Spec) field.Field { return field.NewBitmap(spec) },
	}

	prefixesExtToInt = map[string]prefix.Prefixer{
		"ASCII.Fixed":  prefix.ASCII.Fixed,
		"ASCII.L":      prefix.ASCII.L,
		"ASCII.LL":     prefix.ASCII.LL,
		"ASCII.LLL":    prefix.ASCII.LLL,
		"ASCII.LLLL":   prefix.ASCII.LLLL,
		"BCD.Fixed":    prefix.BCD.Fixed,
		"BCD.L":        prefix.BCD.L,
		"BCD.LL":       prefix.BCD.LL,
		"BCD.LLL":      prefix.BCD.LLL,
		"BCD.LLLL":     prefix.BCD.LLLL,
		"Hex.Fixed":    prefix.Hex.Fixed,
		"Hex.L":        prefix.Hex.L,
		"Hex.LL":       prefix.Hex.LL,
		"Hex.LLL":      prefix.Hex.LLL,
		"Hex.LLLL":     prefix.Hex.LLLL,
		"EBCDIC.Fixed": prefix.EBCDIC.Fixed,
		"EBCDIC.L":     prefix.EBCDIC.L,
		"EBCDIC.LL":    prefix.EBCDIC.LL,
		"EBCDIC.LLL":   prefix.EBCDIC.LLL,
		"EBCDIC.LLLL":  prefix.EBCDIC.LLLL,
		"Binary.Fixed": prefix.Binary.Fixed,
	}

	encodingsExtToInt = map[string]encoding.Encoder{
		"ASCII":  encoding.ASCII,
		"BCD":    encoding.BCD,
		"EBCDIC": encoding.EBCDIC,
		"Binary": encoding.Binary,
		"Hex":    encoding.Hex,
		"LBCD":   encoding.LBCD,
	}

	encodingsIntToExt = map[string]string{
		"asciiEncoder":  "ASCII",
		"bcdEncoder":    "BCD",
		"ebcdicEncoder": "EBCDIC",
		"binaryEncoder": "Binary",
		"hexEncoder":    "Hex",
		"lBCDEncoder":   "LBCD",
	}

	paddersIntToExt = map[string]string{
		"leftPadder": "Left",
		"nonePadder": "None",
	}

	paddersExtToInt = map[string]func(pad string) padding.Padder{
		"Left": func(pad string) padding.Padder {
			if runes := []rune(pad); len(runes) == 1 {
				return padding.Left(runes[0])
			}
			return nil
		},
		"None": func(pad string) padding.Padder { return padding.None },
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
	Pad  string `json:"pad" xml:"pad"`
}

func (builder *messageSpecBuilder) ImportJSON(raw []byte) (*iso8583.MessageSpec, error) {
	dummySpec := specDummy{}
	err := json.Unmarshal(raw, &dummySpec)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling spec: %v", err)
	}

	if len(dummySpec.Fields) == 0 {
		return nil, fmt.Errorf("no fields in JSON file")
	}

	spec := iso8583.MessageSpec{
		Name:   dummySpec.Name,
		Fields: make(map[int]field.Field),
	}

	for key, dummyField := range dummySpec.Fields {
		index := 0
		_, err := fmt.Sscanf(key, FieldPrefix+"%d", &index)
		if err != nil {
			return nil, fmt.Errorf("invalid field index, index's format is `Field`+index ")
		}

		constructor := fieldConstructor[dummyField.Type]
		if constructor == nil {
			return nil, fmt.Errorf("no constructor for filed type: %s for field: %d", dummyField.Type, index)
		}

		enc := encodingsExtToInt[dummyField.Enc]
		if enc == nil {
			return nil, fmt.Errorf("unknown encoding: %s for field: %d", dummyField.Enc, index)
		}

		pref := prefixesExtToInt[dummyField.Prefix]
		if pref == nil {
			return nil, fmt.Errorf("unknown prefix: %s for field: %d", dummyField.Prefix, index)
		}

		var padder padding.Padder
		if dummyField.Padding != nil {
			if padderConstructor := paddersExtToInt[dummyField.Padding.Type]; padderConstructor != nil {
				padder = padderConstructor(dummyField.Padding.Pad)
			}
		}

		spec.Fields[index] = constructor(&field.Spec{
			Length:      dummyField.Length,
			IDLength:    dummyField.IDLength,
			Description: dummyField.Description,
			Enc:         enc,
			Pref:        pref,
			Pad:         padder,
		})

	}

	return &spec, nil
}

func (builder *messageSpecBuilder) ExportJSON(origSpec *iso8583.MessageSpec) ([]byte, error) {
	if origSpec == nil {
		return nil, fmt.Errorf("invalid message spec")
	}

	dummyJson := specDummy{
		Name:   origSpec.Name,
		Fields: make(map[string]fieldDummy),
	}

	for index, origField := range origSpec.Fields {
		dummyField := fieldDummy{}

		spec := origField.Spec()

		dummyField.Length = spec.Length
		dummyField.IDLength = spec.IDLength
		dummyField.Description = spec.Description

		if spec.Enc == nil {
			return nil, fmt.Errorf("missing required spec Enc for field %d", index)
		}

		if spec.Pref == nil {
			return nil, fmt.Errorf("missing required spec Pref for field %d", index)
		}

		// set encoding
		encType := reflect.TypeOf(spec.Enc).Elem().Name()
		if enc, found := encodingsIntToExt[encType]; found {
			dummyField.Enc = enc
		} else {
			return nil, fmt.Errorf("unknown encoding type: %s", encType)
		}

		// set prefixer
		// Inspect() makes a trick for us when we don't have to implement
		// any internal to external representation logic.
		dummyField.Prefix = spec.Pref.Inspect()

		// set padding
		if spec.Pad != nil {
			paddingType := reflect.TypeOf(spec.Pad).Elem().Name()
			if padder, found := paddersIntToExt[paddingType]; found {
				pad := spec.Pad.Inspect()
				dummyPadding := padDummy{
					Type: padder,
					Pad:  string(pad),
				}
				dummyField.Padding = &dummyPadding
			} else {
				return nil, fmt.Errorf("unknown padding type: %s", paddingType)
			}
		}

		fieldType := reflect.TypeOf(origField).Elem().Name()
		fieldName := fmt.Sprintf(FieldPrefix+"%03d", index)

		dummyField.Type = fieldType

		dummyJson.Fields[fieldName] = dummyField
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
