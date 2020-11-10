// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"

	"github.com/Intermernet/ebcdic"
	"github.com/moov-io/iso8583/pkg/utils"
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

func (e *Element) Bytes() ([]byte, error) {
	cat := utils.AvailableTypeCategory[e.Type]
	switch cat {
	case utils.EncodingCatCharacter:
		return e.characterEncoding()
	case utils.EncodingCatBinary:
		return e.binaryEncoding()
	case utils.EncodingCatNumber:
		return e.numberEncoding()
	}
	return nil, errors.New(utils.ErrInvalidEncoder)
}

func (e *Element) Load(raw []byte) (int, error) {
	cat := utils.AvailableTypeCategory[e.Type]
	switch cat {
	case utils.EncodingCatCharacter:
		return e.characterDecoding(raw)
	case utils.EncodingCatBinary:
		return e.binaryDecoding(raw)
	case utils.EncodingCatNumber:
		return e.numberDecoding(raw)
	}
	return 0, errors.New(utils.ErrInvalidEncoder)
}

func (e *Element) characterEncoding() ([]byte, error) {
	var value []byte
	var err error

	if e.Encoding == utils.EncodingAscii {
		value, err = utils.UTF8ToWindows1252(e.Value)
	} else if e.Encoding == utils.EncodingEbcdic {
		value = ebcdic.Encode(e.Value)
	} else {
		return nil, errors.New(utils.ErrInvalidEncoder)
	}

	if err != nil {
		return nil, err
	}
	if len(value) > e.Length {
		return nil, fmt.Errorf(utils.ErrValueTooLong, "character", e.Length, len(value))
	}
	if e.Fixed {
		return value, nil
	}

	lenStr := fmt.Sprintf("%02d", len(value))
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch e.LengthEncoding {
	case utils.EncodingAscii:
		lenVal = contentLen
		if len(lenVal) > 2 {
			return nil, errors.New(utils.ErrInvalidLengthHead)
		}
	case utils.EncodingRBcd:
		lenVal, err = utils.RBcd(contentLen)
		if err != nil {
			return nil, err
		}
		if len(lenVal) > 1 {
			return nil, errors.New(utils.ErrInvalidLengthHead)
		}
	case utils.EncodingBcd:
		lenVal, err = utils.Bcd(contentLen)
		if err != nil {
			return nil, err
		}
		if len(lenVal) > 1 {
			return nil, errors.New(utils.ErrInvalidLengthHead)
		}
	default:
		return nil, errors.New(utils.ErrInvalidLengthEncoder)
	}
	return append(lenVal, value...), nil
}

func (e *Element) numberEncoding() ([]byte, error) {
	return nil, nil
}

func (e *Element) binaryEncoding() ([]byte, error) {
	return nil, nil
}

func (e *Element) characterDecoding(raw []byte) (int, error) {
	var value []byte
	var err error
	var contentLen int
	var read int

	if !e.Fixed {
		lenSize := len(strconv.Itoa(e.Length))
		bcdSize := lenSize / 2
		if lenSize%2 != 0 {
			bcdSize++
		}

		switch e.LengthEncoding {
		case utils.EncodingAscii:
			contentLen, err = strconv.Atoi(string(raw[:lenSize]))
			if err != nil {
				return 0, errors.New(utils.ErrParseLengthFailed + ": " + string(raw[:lenSize]))
			}
			read = lenSize
		case utils.EncodingRBcd:
			lenVal, err := utils.RBcdAscii(raw[:bcdSize], lenSize)
			if err != nil {
				return 0, err
			}
			contentLen, err = strconv.Atoi(string(lenVal))
			if err != nil {
				return 0, errors.New(utils.ErrParseLengthFailed + ": " + string(raw[:lenSize]))
			}
			read = bcdSize
		case utils.EncodingBcd:
			lenVal, err := utils.BcdAscii(raw[:bcdSize], lenSize)
			if err != nil {
				return 0, err
			}
			contentLen, err = strconv.Atoi(string(lenVal))
			if err != nil {
				return 0, errors.New(utils.ErrParseLengthFailed + ": " + string(raw[:lenSize]))
			}
			read = bcdSize
		default:
			return 0, errors.New(utils.ErrInvalidLengthEncoder)
		}
	}

	if contentLen == 0 {
		contentLen = e.Length
	} else {
		e.DataLength = contentLen
	}

	if e.Encoding == utils.EncodingAscii {
		value, err = utils.UTF8ToWindows1252(raw[read : read+contentLen])
	} else if e.Encoding == utils.EncodingEbcdic {
		value = ebcdic.Decode(raw[read : read+contentLen])
	} else {
		return 0, errors.New(utils.ErrInvalidEncoder)
	}

	e.Value = make([]byte, len(value))
	copy(e.Value, value)

	return read, nil
}

func (e *Element) numberDecoding(raw []byte) (int, error) {
	return 0, nil
}

func (e *Element) binaryDecoding(raw []byte) (int, error) {
	return 0, nil
}
