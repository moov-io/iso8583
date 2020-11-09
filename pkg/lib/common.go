// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/moov-io/iso8583/pkg/utils"
	"strconv"
)

// data element, CommonType + Value
type Element struct {
	Type           string
	Length         int
	Format         string
	Encoding       string
	Fixed          bool
	LengthEncoding string

	DataLength int
	Value      []byte
}

func (e *Element) SetType(_type *utils.ElementType) {
	e.Type = _type.Type
	e.Length = _type.Length
	e.Format = _type.Format
	e.Encoding = _type.Encoding
	e.Fixed = _type.Fixed
	e.LengthEncoding = _type.LengthEncoding
}

func (e *Element) Validate() error {
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
	if e.Type == utils.ElementTypeNumeric {
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

func (e *Element) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if e.Type == utils.ElementTypeNumeric {
		ret, err := strconv.Atoi(string(e.Value))
		if err != nil {
			return err
		}
		return encoder.EncodeElement(ret, start)
	}
	s := fmt.Sprintf("%s", e.Value)
	return encoder.EncodeElement(s, start)
}
