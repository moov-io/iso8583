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
	"strings"

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
	DataLength     int
	Value          []byte // raw data without any encoding, equal size of value and length (data length) of element
}

func (e *Element) SetType(_type *utils.ElementType) {
	e.Type = _type.Type
	e.Length = _type.Length
	e.Format = _type.Format
	e.Encoding = _type.Encoding
	e.Fixed = _type.Fixed
	e.LengthEncoding = _type.LengthEncoding
	e.extendBinaryData()
}

func (e *Element) Validate() error {
	return nil
}

func (e *Element) extendBinaryData() {
	cat := utils.AvailableTypeCategory[e.Type]
	if cat == utils.EncodingCatBinary && len(e.Value) < e.Length {
		newData := fmt.Sprintf("%-"+strconv.Itoa(e.Length)+"s", string(e.Value))
		newData = strings.ReplaceAll(newData, " ", "0")
		e.Value = make([]byte, e.Length)
		copy(e.Value, newData)
	}
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
	e.extendBinaryData()
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
	dataLen := e.Length
	if !e.Fixed {
		dataLen = e.DataLength
	}
	if dataLen > len(e.Value) {
		return nil, fmt.Errorf(utils.ErrValueTooLong, "byte", dataLen, len(e.Value))
	}

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

	if e.Fixed {
		return value, nil
	}

	lenEncode, err := e.lengthEncoding(value)
	if err != nil {
		return nil, err
	}

	return append(lenEncode, value...), nil
}

func (e *Element) numberEncoding() ([]byte, error) {
	var value []byte
	var err error

	if e.Encoding == utils.EncodingChar {
		value = e.Value
	} else if e.Encoding == utils.EncodingBcd {
		value, err = utils.Bcd(e.Value)
	} else if e.Encoding == utils.EncodingRBcd {
		value, err = utils.RBcd(e.Value)
	} else {
		return nil, errors.New(utils.ErrInvalidEncoder)
	}
	if err != nil {
		return nil, err
	}

	if e.Fixed {
		return value, nil
	}

	lenEncode, err := e.lengthEncoding(value)
	if err != nil {
		return nil, err
	}

	return append(lenEncode, value...), nil
}

func (e *Element) binaryEncoding() ([]byte, error) {
	var value []byte

	if e.Length != len(e.Value) {
		return nil, errors.New(utils.ErrBadBinary)
	}

	if e.Encoding == utils.EncodingChar {
		value = e.Value
	} else if e.Encoding == utils.EncodingHex {
		bitNum, err := strconv.ParseUint(string(e.Value), 2, e.Length)
		if err != nil {
			return nil, err
		}
		hexStr := fmt.Sprintf("%0"+strconv.Itoa(e.Length/4)+"s", strconv.FormatUint(bitNum, 16))
		value = []byte(hexStr)
	} else {
		return nil, errors.New(utils.ErrInvalidEncoder)
	}

	return value, nil
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
		if err != nil {
			return 0, err
		}
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

func (e *Element) lengthEncoding(value []byte) ([]byte, error) {
	contentLen := []byte(fmt.Sprintf("%02d", len(value)))

	var encode []byte
	var err error

	switch e.LengthEncoding {
	case utils.EncodingAscii:
		encode = contentLen
		if len(encode) > 2 {
			return nil, errors.New(utils.ErrInvalidLengthHead)
		}
	case utils.EncodingRBcd:
		encode, err = utils.RBcd(contentLen)
		if err != nil {
			return nil, err
		}
		if len(encode) > 1 {
			return nil, errors.New(utils.ErrInvalidLengthHead)
		}
	case utils.EncodingBcd:
		encode, err = utils.Bcd(contentLen)
		if err != nil {
			return nil, err
		}
		if len(encode) > 1 {
			return nil, errors.New(utils.ErrInvalidLengthHead)
		}
	default:
		return nil, errors.New(utils.ErrInvalidLengthEncoder)
	}

	return encode, nil
}
