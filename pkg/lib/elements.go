// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"sort"
	"strconv"

	"github.com/moov-io/iso8583/pkg/utils"
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
func NewDataElements(spec *utils.Specification) (*dataElements, error) {
	if spec == nil || spec.Elements == nil || spec.Encoding == nil {
		return nil, errors.New(utils.ErrInvalidSpecification)
	}
	return &dataElements{
		elements: make(map[int]*Element),
		spec:     spec,
	}, nil
}

// data elements of the iso8583 message
type dataElements struct {
	elements map[int]*Element
	spec     *utils.Specification
}

// Validate check validation of field
func (e *dataElements) Validate() error {
	for _, _element := range e.elements {
		if err := _element.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Bytes encode field to bytes
func (e *dataElements) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	for _, key := range e.Keys() {
		element, exist := e.elements[key]
		if exist {
			value, err := element.Bytes()
			if err != nil {
				return nil, err
			}
			buf.Write(value)
		}
	}

	return buf.Bytes(), nil
}

// Load decode field from bytes
func (e *dataElements) Load(raw []byte) (int, error) {
	return 0, nil
}

func (e *dataElements) Keys() []int {
	var keys []int
	for k := range e.elements {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// Customize unmarshal of json
func (e *dataElements) UnmarshalJSON(b []byte) error {
	if e.spec == nil {
		return errors.New(utils.ErrNonExistSpecification)
	}
	if err := json.Unmarshal(b, &e.elements); err != nil {
		return err
	}
	for key, elm := range e.elements {
		spec, err := e.spec.Elements.Get(key)
		if err != nil {
			return err
		}
		_type, err := spec.Parse()
		if err != nil {
			return err
		}
		_type.SetEncoding(e.spec.Encoding)
		elm.setType(_type)
	}
	return nil
}

// Customize marshal of json
func (e *dataElements) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, key := range e.Keys() {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.Write([]byte(`"` + strconv.Itoa(key) + `"`))
		buf.WriteString(":")
		val, err := json.Marshal(e.elements[key])
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}

// Customize unmarshal of xml
func (e *dataElements) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}

	for _, key := range e.Keys() {
		t := xml.StartElement{
			Name: xml.Name{Local: utils.DataElementXmlName},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: utils.DataElementAttrNumber}, Value: strconv.Itoa(key)},
			},
		}
		tokens = append(tokens, t, xml.CharData(e.elements[key].String()), xml.EndElement{Name: t.Name})
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

// Customize marshal of xml
func (e *dataElements) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	if e.spec == nil || e.spec.Encoding == nil {
		return errors.New(utils.ErrNonExistSpecification)
	}

	var dummy xmlDataElement
	err := decoder.DecodeElement(&dummy, &start)
	if err != nil {
		return err
	}

	for _, element := range dummy.Elements {
		var dataElement Element

		spec, err := e.spec.Elements.Get(element.Number)
		if err != nil {
			return err
		}
		_type, err := spec.Parse()
		if err != nil {
			return err
		}

		_type.SetEncoding(e.spec.Encoding)
		dataElement.setType(_type)
		dataElement.Value = make([]byte, len(element.Text))
		copy(dataElement.Value, element.Text)
		e.elements[element.Number] = &dataElement
	}

	return nil
}
