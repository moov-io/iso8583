package iso8583

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"unicode/utf8"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

func newMessageSpecMarshaller() *messageSpecMarshaller {
	return &messageSpecMarshaller{
		XMLName:    xml.Name{Local: "Specification"},
		Fields:     make(map[string]specMarshaller, 0),
		FieldArray: make([]specMarshaller, 0),
	}
}

type paddingMarshaller struct {
	Type   string  `xml:"Type" json:"type"`
	Padder *string `xml:"Padder" json:"padder,omitempty"`
}

type specMarshaller struct {
	XMLName     xml.Name           `json:"-"`
	Type        string             `xml:"Type" json:"type"`
	Length      int                `xml:"Length" json:"length"`
	Encoding    string             `xml:"Enc" json:"enc"`
	Pref        string             `xml:"Pref" json:"pref"`
	Description *string            `xml:"Description,omitempty" json:"description,omitempty"`
	Identifier  *string            `xml:"Identifier,omitempty" json:"identifier,omitempty"`
	Pad         *paddingMarshaller `xml:"Pad,omitempty" json:"pad,omitempty"`
}

func (s *specMarshaller) setType(fieldElement field.Field) {

	typeConvertTable := map[string]string{
		"field.String":  "String",
		"field.Bitmap":  "Bitmap",
		"field.Numeric": "Numeric",
	}

	s.Type = typeConvertTable[reflect.TypeOf(fieldElement).Elem().String()]
}

func (s *specMarshaller) setEncoding(spec *field.Spec) {

	encConvertTable := map[string]string{
		"encoding.asciiEncoder":  "ASCII",
		"encoding.hexEncoder":    "Hex",
		"encoding.bcdEncoder":    "BCD",
		"encoding.binaryEncoder": "Binary",
		"encoding.lBCDEncoder":   "LBCD",
	}

	if spec.Enc == nil {
		return
	}

	s.Encoding = encConvertTable[reflect.TypeOf(spec.Enc).Elem().String()]
}

func (s *specMarshaller) setPrefix(spec *field.Spec) {

	prefConvertTable := map[string]string{
		"prefix.asciiFixedPrefixer.Fixed":  "ASCII.Fixed",
		"prefix.asciiVarPrefixer.L":        "ASCII.L",
		"prefix.asciiVarPrefixer.LL":       "ASCII.LL",
		"prefix.asciiVarPrefixer.LLL":      "ASCII.LLL",
		"prefix.asciiVarPrefixer.LLLL":     "ASCII.LLLL",
		"prefix.bcdFixedPrefixer.Fixed":    "BCD.Fixed",
		"prefix.bcdVarPrefixer.L":          "BCD.L",
		"prefix.bcdVarPrefixer.LL":         "BCD.LL",
		"prefix.bcdVarPrefixer.LLL":        "BCD.LLL",
		"prefix.bcdVarPrefixer.LLLL":       "BCD.LLLL",
		"prefix.binaryFixedPrefixer.Fixed": "Binary.Fixed",
		"prefix.hexFixedPrefixer.Fixed":    "Hex.Fixed",
	}

	if spec.Enc == nil {
		return
	}

	s.Pref = prefConvertTable[reflect.TypeOf(spec.Pref).Elem().String()+"."+spec.Pref.InspectName()]
}

func (s *specMarshaller) setPadding(spec *field.Spec) {
	padConvertTable := map[string]string{
		"padding.leftPadder": "Left",
		"padding.nonePadder": "None",
	}

	if spec.Pad == nil {
		return
	}

	s.Pad = &paddingMarshaller{}
	s.Pad.Type = padConvertTable[reflect.TypeOf(spec.Pad).Elem().String()]
	if spec.Pad.Inspect() != nil {
		s.Pad.Padder = spec.Pad.Inspect()
	}
}

func (s *specMarshaller) getEncoder(encoderName string) encoding.Encoder {
	encoderTable := map[string]encoding.Encoder{
		"ASCII":  encoding.ASCII,
		"Hex":    encoding.Hex,
		"BCD":    encoding.BCD,
		"Binary": encoding.Binary,
		"LBCD":   encoding.LBCD,
	}

	encoder, ok := encoderTable[encoderName]

	if !ok {
		return nil
	}

	return encoder
}

func (s *specMarshaller) getPrefixer(encoderName string) prefix.Prefixer {
	prefixerTable := map[string]prefix.Prefixer{
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
		"Binary.Fixed": prefix.Binary.Fixed,
		"Hex.Fixed":    prefix.Hex.Fixed,
	}

	prefixer, ok := prefixerTable[encoderName]

	if !ok {
		return nil
	}

	return prefixer
}

func (s *specMarshaller) getPadder(pad *paddingMarshaller) padding.Padder {

	switch pad.Type {
	case "Left":
		padLetter := '0'
		if pad.Padder != nil {
			r, size := utf8.DecodeRune([]byte(*pad.Padder))
			if size > 0 {
				padLetter = r
			}
		}
		return padding.Left(padLetter)
	case "None":
		return padding.None
	}

	return nil
}

func (s *specMarshaller) createSpecificationField() field.Field {

	spec := &field.Spec{
		Length: s.Length,
		Enc:    s.getEncoder(s.Encoding),
		Pref:   s.getPrefixer(s.Pref),
	}
	if s.Description != nil {
		spec.Description = *s.Description
	}
	if s.Identifier != nil {
		spec.Identifier = *s.Identifier
	}
	if s.Pad != nil {
		spec.Pad = s.getPadder(s.Pad)
	}

	switch s.Type {
	case "String":
		return field.NewString(spec)
	case "Bitmap":
		return field.NewBitmap(spec)
	case "Numeric":
		return field.NewNumeric(spec)
	}

	return nil
}

type messageSpecMarshaller struct {
	XMLName    xml.Name                  `xml:"Specification" json:"-"`
	Fields     map[string]specMarshaller `xml:"-"`
	FieldArray []specMarshaller          `json:"-"`
}

func (m *messageSpecMarshaller) addSpecMarshaller(index int, fieldElement field.Field) error {

	stringToPointer := func(value string) *string {
		if len(value) == 0 {
			return nil
		}
		return &value
	}

	spec := fieldElement.Spec()
	if spec == nil {
		return fmt.Errorf("there is not specification of %d element", index)
	}

	newDummySpec := specMarshaller{
		Description: stringToPointer(spec.Description),
		Identifier:  stringToPointer(spec.Identifier),
		Length:      spec.Length,
	}

	newDummySpec.setType(fieldElement)
	newDummySpec.setEncoding(spec)
	newDummySpec.setPrefix(spec)
	newDummySpec.setPadding(spec)

	fieldName := fmt.Sprintf("%03d", index)

	newDummySpec.XMLName = xml.Name{Local: fieldName}
	m.Fields[fieldName] = newDummySpec

	newDummySpec.XMLName = xml.Name{Local: "F" + fieldName}
	m.FieldArray = append(m.FieldArray, newDummySpec) //  should not start with numbers for xml tag name

	return nil
}
