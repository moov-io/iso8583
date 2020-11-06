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

// dummy struct for xml un-marshaling
type xmlDataElement struct {
	XMLName  xml.Name `xml:"DataElements"`
	Text     string   `xml:",chardata"`
	Elements []struct {
		Text   string `xml:",chardata"`
		Number int    `xml:"Number,attr"`
	} `xml:"Element"`
}

// create data elements of message with specification
func NewDataElements(spec Specification) (*DataElements, error) {
	if spec.Elements == nil || spec.Encoding == nil {
		return nil, errors.New("has invalid specification")
	}
	return &DataElements{
		Elements: make(map[int]*Element),
		Spec:     spec,
	}, nil
}

// data element, CommonType + Value
type Element struct {
	Type  CommonType
	Value []byte
}

func (e *Element) SetType(_type *CommonType) {
	e.Type = *_type
}

func (e *Element) Validate() error {
	err := e.Type.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (e *Element) String() string {
	return fmt.Sprintf("%s", e.Value)
}

func (e *Element) UnmarshalJSON(b []byte) error {
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

func (e *Element) MarshalJSON() ([]byte, error) {
	if e.Type.Type == ElementTypeNumeric {
		ret, err := strconv.Atoi(string(e.Value))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return json.Marshal(ret)
	}
	return json.Marshal(fmt.Sprintf("%s", e.Value))
}

func (e *Element) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := decoder.DecodeElement(&s, &start); err != nil {
		return err
	}
	return nil
}

// data elements of the iso8583 message
type DataElements struct {
	Elements map[int]*Element
	Spec     Specification
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
	var convert map[int]*Element
	convert = e.Elements
	if err := json.Unmarshal(b, &convert); err != nil {
		return err
	}
	for key, elm := range convert {
		spec, err := e.Spec.Elements.Get(key)
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

func (e *DataElements) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var dummy xmlDataElement
	err := decoder.DecodeElement(&dummy, &start)
	if err != nil {
		return err
	}

	for _, element := range dummy.Elements {
		var dataElement Element

		spec, err := e.Spec.Elements.Get(element.Number)
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
