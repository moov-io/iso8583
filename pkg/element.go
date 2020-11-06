// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package pkg

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

// create data elements of message with specification
func NewDataElements(spec Specifications) *DataElements {
	return &DataElements{
		Elements:       make(map[int]*DataElement),
		Specifications: spec,
	}
}

// general element type for all of the data representation attributes
type CommonType struct {
	Type   string
	Length int
	Fixed  bool
	Format string
}

func (e CommonType) DataElement(value []byte) (*DataElement, error) {
	if value == nil {
		return nil, errors.New("invalid element value")
	}
	_new := DataElement{
		Type:   e.Type,
		Length: e.Length,
		Fixed:  e.Fixed,
		Format: e.Format,
		Value:  make([]byte, len(value)),
	}
	copy(_new.Value, value)

	return &_new, nil
}

// data element, CommonType + Value
type DataElement struct {
	Type   string
	Length int
	Fixed  bool
	Format string
	Value  []byte
}

func (e *DataElement) SetType(_type *CommonType) {
	e.Type = _type.Type
	e.Length = _type.Length
	e.Fixed = _type.Fixed
	e.Format = _type.Format
}

func (e *DataElement) Validate() error {
	return nil
}

func (e *DataElement) String() string {
	return fmt.Sprintf("%s", e.Value)
}

func (e *DataElement) UnmarshalJSON(b []byte) error {
	_, err := strconv.ParseFloat(string(b), 64)
	if err == nil {
		e.Value = make([]byte, len(b))
		copy(e.Value, b)
	} else {
		var value string
		err := json.Unmarshal(b, &value)
		if err != nil {
			return err
		}
		e.Value = make([]byte, len(value))
		copy(e.Value, value)
	}
	return nil
}

func (e *DataElement) MarshalJSON() ([]byte, error) {
	if e.Type == ElementTypeNumeric {
		ret, err := strconv.Atoi(string(e.Value))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return json.Marshal(ret)
	}
	return json.Marshal(fmt.Sprintf("%s", e.Value))
}

func (e *DataElement) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := decoder.DecodeElement(&s, &start); err != nil {
		return err
	}
	return nil
}

// data elements of the iso8583 message
type DataElements struct {
	Elements       map[int]*DataElement
	Specifications Specifications
}

func (e *DataElements) Validate() error {
	for _, _element := range e.Elements {
		if err := _element.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e *DataElements) UnmarshalJSON(b []byte) error {
	var convert map[int]*DataElement
	convert = e.Elements
	if err := json.Unmarshal(b, &convert); err != nil {
		return err
	}
	for key, elm := range convert {
		spec, err := e.Specifications.Get(key)
		if err != nil {
			return errors.New("don't exist specification")
		}
		_type, err := spec.ElementType()
		if err != nil {
			return err
		}
		elm.SetType(_type)
	}
	return nil
}

func (e *DataElements) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Elements)
}

func (e *DataElements) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}
	for key, value := range e.Elements {
		t := xml.StartElement{
			Name: xml.Name{Local: DataElementXmlName},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: DataElementAttrNumber}, Value: strconv.Itoa(key)},
			},
		}
		tokens = append(tokens, t, xml.CharData(value.String()), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{Name: start.Name})
	for _, t := range tokens {
		err := encoder.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	return encoder.Flush()
}

// dummy struct for xml un-marshaling
type xmlDataElement struct {
	XMLName  xml.Name `xml:"DataElements"`
	Text     string   `xml:",chardata"`
	Elements []struct {
		Text   string `xml:",chardata"`
		Number int    `xml:"Number,attr"`
	} `xml:"Element"`
}

func (e *DataElements) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var dummy xmlDataElement
	err := decoder.DecodeElement(&dummy, &start)
	if err != nil {
		return err
	}

	for _, element := range dummy.Elements {
		var dataElement DataElement

		spec, err := e.Specifications.Get(element.Number)
		if err != nil {
			return err
		}
		_type, err := spec.ElementType()
		if err != nil {
			return err
		}

		dataElement.SetType(_type)
		dataElement.Value = make([]byte, len(element.Text))
		copy(dataElement.Value, element.Text)
		e.Elements[element.Number] = &dataElement
	}

	return nil
}
