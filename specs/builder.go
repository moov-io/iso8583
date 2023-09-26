package specs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strconv"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	moovsort "github.com/moov-io/iso8583/sort"
	"github.com/moov-io/iso8583/utils"
)

type FieldConstructorFunc func(spec *field.Spec) field.Field

var (
	FieldConstructor = map[string]FieldConstructorFunc{
		"String":    func(spec *field.Spec) field.Field { return field.NewString(spec) },
		"Numeric":   func(spec *field.Spec) field.Field { return field.NewNumeric(spec) },
		"Binary":    func(spec *field.Spec) field.Field { return field.NewBinary(spec) },
		"Bitmap":    func(spec *field.Spec) field.Field { return field.NewBitmap(spec) },
		"Composite": func(spec *field.Spec) field.Field { return field.NewComposite(spec) },
	}

	PrefixesExtToInt = map[string]prefix.Prefixer{
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
		"Binary.L":     prefix.Binary.L,
		"Binary.LL":    prefix.Binary.LL,
		"Binary.LLL":   prefix.Binary.LLL,
		"Binary.LLLL":  prefix.Binary.LLLL,
		"BerTLV":       prefix.BerTLV,
	}

	EncodingsExtToInt = map[string]encoding.Encoder{
		"ASCII":      encoding.ASCII,
		"BCD":        encoding.BCD,
		"EBCDIC":     encoding.EBCDIC,
		"Binary":     encoding.Binary,
		"HexToASCII": encoding.BytesToASCIIHex,
		"ASCIIToHex": encoding.ASCIIHexToBytes,
		"LBCD":       encoding.LBCD,
		"BerTLVTag":  encoding.BerTLVTag,
	}

	EncodingsIntToExt = map[string]string{
		"asciiEncoder":      "ASCII",
		"bcdEncoder":        "BCD",
		"ebcdicEncoder":     "EBCDIC",
		"binaryEncoder":     "Binary",
		"hexToASCIIEncoder": "HexToASCII",
		"asciiToHexEncoder": "ASCIIToHex",
		"lBCDEncoder":       "LBCD",
	}

	PaddersIntToExt = map[string]string{
		"leftPadder":  "Left",
		"rightPadder": "Right",
		"nonePadder":  "None",
	}

	PaddersExtToInt = map[string]func(pad string) padding.Padder{
		"Left": func(pad string) padding.Padder {
			if runes := []rune(pad); len(runes) == 1 {
				return padding.Left(runes[0])
			}
			return nil
		},
		"Right": func(pad string) padding.Padder {
			if runes := []rune(pad); len(runes) == 1 {
				return padding.Right(runes[0])
			}
			return nil
		},
		"None": func(pad string) padding.Padder { return padding.None },
	}

	SortExtToInt = map[string]moovsort.StringSlice{
		"StringsByInt": moovsort.StringsByInt,
		"StringsByHex": moovsort.StringsByHex,
	}
)

var Builder MessageSpecBuilder = &messageSpecBuilder{}

type MessageSpecBuilder interface {
	ImportJSON([]byte) (*iso8583.MessageSpec, error)
	ExportJSON(spec *iso8583.MessageSpec) ([]byte, error)
}

type messageSpecBuilder struct{}

type specDummy struct {
	Name   string          `json:"name,omitempty" xml:"name,omitempty"`
	Fields orderedFieldMap `json:"fields,omitempty" xml:"fields,omitempty"`
}

type fieldDummy struct {
	Type        string                 `json:"type,omitempty" xml:"type,omitempty"`
	Length      int                    `json:"length,omitempty" xml:"length,omitempty"`
	Description string                 `json:"description,omitempty" xml:"description,omitempty"`
	Enc         string                 `json:"enc,omitempty" xml:"enc,omitempty"`
	Prefix      string                 `json:"prefix,omitempty" xml:"prefix,omitempty"`
	Padding     *paddingDummy          `json:"padding,omitempty" xml:"padding,omitempty"`
	Tag         *tagDummy              `json:"tag,omitempty" xml:"tag,omitempty"`
	Subfields   map[string]*fieldDummy `json:"subfields,omitempty" xml:"subfields:omitempty"`
}

type paddingDummy struct {
	Type string `json:"type" xml:"type"`
	Pad  string `json:"pad" xml:"pad"`
}

type tagDummy struct {
	Length  int           `json:"length,omitempty" xml:"length,omitempty"`
	Enc     string        `json:"enc,omitempty" xml:"enc,omitempty"`
	Padding *paddingDummy `json:"padding,omitempty" xml:"padding,omitempty"`
	Sort    string        `json:"sort,omitempty" xml:"sort,omitempty"`
}

func importField(dummyField *fieldDummy, index string) (*field.Spec, error) {
	fieldSpec := &field.Spec{
		Length:      dummyField.Length,
		Description: dummyField.Description,
	}

	fieldSpec.Pref = PrefixesExtToInt[dummyField.Prefix]
	if fieldSpec.Pref == nil {
		return nil, fmt.Errorf("unknown prefix: %s for field: %s", dummyField.Prefix, index)
	}

	if dummyField.Padding != nil {
		if padderConstructor := PaddersExtToInt[dummyField.Padding.Type]; padderConstructor != nil {
			fieldSpec.Pad = padderConstructor(dummyField.Padding.Pad)
		}
	}

	if len(dummyField.Subfields) == 0 {
		fieldSpec.Enc = EncodingsExtToInt[dummyField.Enc]
		if fieldSpec.Enc == nil {
			return nil, fmt.Errorf("unknown encoding: %s for field: %s", dummyField.Enc, index)
		}
	} else {
		fieldSpec.Subfields = map[string]field.Field{}
		for key, dummyField := range dummyField.Subfields {
			subfieldSpec, err := importField(dummyField, key)
			if err != nil {
				return nil, err
			}
			constructor := FieldConstructor[dummyField.Type]
			if constructor == nil {
				return nil, fmt.Errorf("no constructor for filed type: %s for field: %s", dummyField.Type, index)
			}
			fieldSpec.Subfields[key] = constructor(subfieldSpec)
		}

		fieldSpec.Tag = &field.TagSpec{
			Length: dummyField.Tag.Length,
		}
		fieldSpec.Tag.Enc = EncodingsExtToInt[dummyField.Tag.Enc]
		if dummyField.Tag.Padding != nil {
			if padderConstructor := PaddersExtToInt[dummyField.Tag.Padding.Type]; padderConstructor != nil {
				fieldSpec.Tag.Pad = padderConstructor(dummyField.Tag.Padding.Pad)
			}
		}
		fieldSpec.Tag.Sort = SortExtToInt[dummyField.Tag.Sort]
	}
	return fieldSpec, nil
}

func (builder *messageSpecBuilder) ImportJSON(raw []byte) (*iso8583.MessageSpec, error) {
	dummySpec := specDummy{}
	err := json.Unmarshal(raw, &dummySpec)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON unmarshal bytes to MessageSpec")
	}

	if len(dummySpec.Fields) == 0 {
		return nil, fmt.Errorf("no fields in JSON file")
	}

	spec := iso8583.MessageSpec{
		Name:   dummySpec.Name,
		Fields: make(map[int]field.Field),
	}

	for key, dummyField := range dummySpec.Fields {
		index, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid field index: %w", err)
		}
		fieldSpec, err := importField(dummyField, key)
		if err != nil {
			return nil, fmt.Errorf("error importing field: %d. %w", index, err)
		}
		constructor := FieldConstructor[dummyField.Type]
		if constructor == nil {
			return nil, fmt.Errorf("no constructor for filed type: %s for field: %d", dummyField.Type, index)
		}
		spec.Fields[index] = constructor(fieldSpec)
	}

	return &spec, nil
}

func exportField(internalField field.Field) (*fieldDummy, error) {
	spec := internalField.Spec()
	fieldType := reflect.TypeOf(internalField).Elem().Name()
	dummyField := &fieldDummy{
		Type:        fieldType,
		Length:      spec.Length,
		Description: spec.Description,
	}

	if spec.Pref == nil {
		return nil, fmt.Errorf("missing required spec.Pref")
	}
	// Inspect() makes a trick for us when we don't have to implement
	// any internal to external representation logic.
	dummyField.Prefix = spec.Pref.Inspect()

	if spec.Pad != nil {
		dummyPad, err := exportPad(spec.Pad)
		if err != nil {
			return nil, err
		}
		dummyField.Padding = dummyPad
	}

	if len(spec.Subfields) == 0 {
		// Encoding only applies to primitive field types
		if spec.Enc == nil {
			return nil, fmt.Errorf("missing required spec.Enc")
		}
		enc, err := exportEnc(spec.Enc)
		if err != nil {
			return nil, err
		}
		dummyField.Enc = enc

	} else {
		dummyField.Subfields = map[string]*fieldDummy{}
		for index, origField := range spec.Subfields {
			f, err := exportField(origField)
			if err != nil {
				return nil, err
			}
			dummyField.Subfields[index] = f
		}

		if spec.Tag != nil {
			tag, err := exportTag(spec.Tag)
			if err != nil {
				return nil, err
			}
			dummyField.Tag = tag
		}
	}

	return dummyField, nil
}

func exportTag(tag *field.TagSpec) (*tagDummy, error) {
	dummy := &tagDummy{
		Length: tag.Length,
	}
	if tag.Pad != nil {
		var err error
		if dummy.Padding, err = exportPad(tag.Pad); err != nil {
			return nil, err
		}
	}

	if tag.Enc != nil {
		var err error
		if dummy.Enc, err = exportEnc(tag.Enc); err != nil {
			return nil, err
		}
	}
	if tag.Sort != nil {
		dummy.Sort = getFunctionName(tag.Sort)
	}
	return dummy, nil

}

func exportPad(pad padding.Padder) (*paddingDummy, error) {
	paddingType := reflect.TypeOf(pad).Elem().Name()
	if padder, found := PaddersIntToExt[paddingType]; found {
		return &paddingDummy{
			Type: padder,
			Pad:  string(pad.Inspect()),
		}, nil
	}
	return nil, fmt.Errorf("unknown padding type: %s", paddingType)
}
func exportEnc(enc encoding.Encoder) (string, error) {
	// set encoding
	encType := reflect.TypeOf(enc).Elem().Name()
	if e, found := EncodingsIntToExt[encType]; found {
		return e, nil
	} else {
		return "", fmt.Errorf("unknown encoding type: %s", encType)
	}
}

func (builder *messageSpecBuilder) ExportJSON(origSpec *iso8583.MessageSpec) ([]byte, error) {
	if origSpec == nil {
		return nil, fmt.Errorf("invalid message spec")
	}
	dummy := specDummy{
		Name:   origSpec.Name,
		Fields: map[string]*fieldDummy{},
	}

	for index, origField := range origSpec.Fields {
		f, err := exportField(origField)
		if err != nil {
			return nil, fmt.Errorf("failed to export field: %d. %w", index, err)
		}
		dummy.Fields[strconv.Itoa(index)] = f
	}

	outputBuffer := new(bytes.Buffer)
	enc := json.NewEncoder(outputBuffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")

	if err := enc.Encode(dummy); err != nil {
		return nil, utils.NewSafeError(err, "failed to perform JSON encoding")
	}

	return outputBuffer.Bytes(), nil
}

type orderedFieldMap map[string]*fieldDummy

func (om orderedFieldMap) MarshalJSON() ([]byte, error) {
	keys := make([]int, 0, len(om))
	for k := range om {
		index, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("converting field index into int: %w", err)
		}
		keys = append(keys, index)
	}

	sort.Ints(keys)

	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, key := range keys {
		strIndex := strconv.Itoa(key)
		b, err := json.Marshal(om[strIndex])
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("\"%s\":", strIndex))
		buf.Write(b)

		// don't add "," if it's the last item
		if i == len(keys)-1 {
			break
		}

		buf.Write([]byte{','})
	}
	buf.Write([]byte{'}'})

	return buf.Bytes(), nil
}

func getFunctionName(foo interface{}) string {
	funcPath := runtime.FuncForPC(reflect.ValueOf(foo).Pointer()).Name()
	return path.Ext(funcPath)[1:]
}
