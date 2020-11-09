// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/moov-io/iso8583/pkg/utils"
	"sort"
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
func NewDataElements(spec *utils.Specification) (*DataElements, error) {
	if spec == nil && spec.Elements == nil || spec.Encoding == nil {
		return nil, errors.New("has invalid specification")
	}
	return &DataElements{
		Elements: make(map[int]*Element),
		Spec:     spec,
	}, nil
}

// data elements of the iso8583 message
type DataElements struct {
	Elements map[int]*Element
	Spec     *utils.Specification `xml:"-" json:"-"`
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
	if e.Spec == nil {
		return errors.New("don't exist specification")
	}
	var convert map[int]*Element
	convert = e.Elements
	if err := json.Unmarshal(b, &convert); err != nil {
		return err
	}
	for key, elm := range convert {
		spec, err := e.Spec.Elements.Get(key)
		if err != nil {
			return err
		}
		_type, err := spec.ElementType()
		if err != nil {
			return err
		}
		_type.SetEncoding(e.Spec.Encoding)
		elm.SetType(_type)
	}
	return nil
}

func (e *DataElements) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	var keys []int

	if e.Elements != nil {
		for k, _ := range e.Elements {
			keys = append(keys, k)
		}
		sort.Ints(keys)
	}

	buf.WriteString("{")
	for i, key := range keys {
		if i != 0 {
			buf.WriteString(",")
		}
		number, err := json.Marshal(strconv.Itoa(key))
		if err != nil {
			return nil, err
		}
		buf.Write(number)
		buf.WriteString(":")
		val, err := json.Marshal(e.Elements[key])
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}

func (e *DataElements) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}
	var keys []int

	if e.Elements != nil {
		for k, _ := range e.Elements {
			keys = append(keys, k)
		}
		sort.Ints(keys)
	}

	for _, key := range keys {
		t := xml.StartElement{
			Name: xml.Name{Local: utils.DataElementXmlName},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: utils.DataElementAttrNumber}, Value: strconv.Itoa(key)},
			},
		}
		tokens = append(tokens, t, xml.CharData(e.Elements[key].String()), xml.EndElement{Name: t.Name})
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
	if e.Spec == nil {
		return errors.New("don't exist specification")
	}

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

		_type.SetEncoding(e.Spec.Encoding)
		dataElement.SetType(_type)
		dataElement.Value = make([]byte, len(element.Text))
		copy(dataElement.Value, element.Text)
		e.Elements[element.Number] = &dataElement
	}

	return nil
}
