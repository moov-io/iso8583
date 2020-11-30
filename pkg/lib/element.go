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

	"github.com/indece-official/go-ebcdic"
	"github.com/moov-io/iso8583/pkg/utils"
)

// data element, CommonType + Value
type Element struct {
	Type           string `xml:"-" json:"-"`
	Length         int    `xml:"-" json:"-"`
	Format         string `xml:"-" json:"-"`
	Encoding       string `xml:"-" json:"-"`
	Fixed          bool   `xml:"-" json:"-"`
	LengthEncoding string `xml:"-" json:"-"`
	DataLength     int    `xml:"-" json:"-"`
	Value          []byte `xml:"-" json:"-"` // raw data without any encoding, equal size of value and length (data length) of element
}

// Validate check validation of field
func (e *Element) Validate() error {
	// Checking available encoding
	err := e.validateWithEncoding()
	if err != nil {
		return err
	}

	// Regex
	err = e.validateWithRegex()
	if err != nil {
		return err
	}

	// Date Format
	err = e.validateWithFormat()
	if err != nil {
		return err
	}

	return nil
}

// String field to string
func (e *Element) String() string {
	return fmt.Sprintf("%s", e.Value)
}

// Bytes encode field to bytes
func (e *Element) Bytes() ([]byte, error) {
	dataLen := e.Length
	if !e.Fixed {
		dataLen = e.DataLength
	}
	if len(e.Value) > dataLen {
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

// Load decode field from bytes
func (e *Element) Load(raw []byte) (int, error) {
	switch utils.AvailableTypeCategory[e.Type] {
	case utils.EncodingCatCharacter:
		return e.characterDecoding(raw)
	case utils.EncodingCatBinary:
		return e.binaryDecoding(raw)
	case utils.EncodingCatNumber:
		return e.numberDecoding(raw)
	}
	return 0, errors.New(utils.ErrInvalidEncoder)
}

// Customize unmarshal of json
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

// Customize marshal of json
func (e *Element) MarshalJSON() ([]byte, error) {
	if e.Type == utils.ElementTypeNumeric {
		ret, err := strconv.Atoi(string(e.Value))
		if err != nil {
			return nil, err
		}
		return json.Marshal(ret)
	}
	return json.Marshal(fmt.Sprintf("%s", e.Value))
}

// Customize unmarshal of xml
func (e *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	return nil
}

// Customize marshal of xml
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

// private functions ...
func (e *Element) characterEncoding() ([]byte, error) {
	var encodingValue []byte
	var err error

	if e.Encoding == utils.EncodingChar {
		encodingValue = e.Value
	} else if e.Encoding == utils.EncodingAscii {
		encodingValue, err = utils.UTF8ToWindows1252(e.Value)
	} else if e.Encoding == utils.EncodingEbcdic {
		encodingValue, err = ebcdic.Encode(string(e.Value), ebcdic.EBCDIC037)
	} else {
		return nil, errors.New(utils.ErrInvalidEncoder)
	}
	if err != nil {
		return nil, err
	}

	paddingValue := e.characterWithPadding(encodingValue)
	if e.Fixed {
		return paddingValue, nil
	}

	lenEncode, err := e.lengthEncoding(encodingValue)
	if err != nil {
		return nil, err
	}

	return append(lenEncode, paddingValue...), nil
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

	paddingValue := e.numberWithPadding(value, true)
	if e.Fixed {
		return paddingValue, nil
	}

	lenEncode, err := e.lengthEncoding(value)
	if err != nil {
		return nil, err
	}

	return append(lenEncode, paddingValue...), nil
}

func (e *Element) binaryEncoding() ([]byte, error) {
	var value []byte

	if e.Length != len(e.Value) {
		return nil, errors.New(utils.ErrBadBinary)
	}

	if e.Encoding == utils.EncodingChar {
		e.extendBinaryData()
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
	read, err := e.lengthDecoding(raw)
	if err != nil {
		return 0, err
	}

	contentLen := e.Length
	if !e.Fixed {
		contentLen = e.DataLength
	}

	var value []byte
	if len(raw) < read+contentLen {
		return 0, errors.New(utils.ErrBadElementData)
	}

	if e.Encoding == utils.EncodingAscii {
		value, err = utils.UTF8ToWindows1252(raw[read : read+contentLen])
	} else if e.Encoding == utils.EncodingEbcdic {
		var str string
		str, err = ebcdic.Decode(raw[read:read+contentLen], ebcdic.EBCDIC037)
		value = []byte(str)
	} else {
		return 0, errors.New(utils.ErrInvalidEncoder)
	}

	if err != nil {
		return 0, err
	}

	e.Value = make([]byte, len(value))
	copy(e.Value, value)
	read += contentLen

	return read, nil
}

func (e *Element) numberDecoding(raw []byte) (int, error) {
	read, err := e.lengthDecoding(raw)
	if err != nil {
		return 0, err
	}

	contentLen := e.Length
	if !e.Fixed {
		contentLen = e.DataLength
	}

	var value []byte
	if e.Encoding == utils.EncodingChar {
		if len(raw) < read+contentLen {
			return 0, errors.New(utils.ErrBadElementData)
		}
		value, err = utils.UTF8ToWindows1252(raw[read : read+contentLen])
		if err != nil {
			return 0, err
		}
	} else if e.Encoding == utils.EncodingRBcd {
		bcdSize := contentLen / 2
		if (contentLen)%2 != 0 {
			bcdSize++
		}
		if len(raw) < read+bcdSize {
			return 0, errors.New(utils.ErrBadElementData)
		}
		value, err = utils.RBcdAscii(raw[read:read+bcdSize], bcdSize)
		if err != nil {
			return 0, err
		}
	} else if e.Encoding == utils.EncodingBcd {
		bcdSize := contentLen / 2
		if (contentLen)%2 != 0 {
			bcdSize++
		}
		if len(raw) < read+bcdSize {
			return 0, errors.New(utils.ErrBadElementData)
		}
		value, err = utils.BcdAscii(raw[read:read+bcdSize], bcdSize)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New(utils.ErrInvalidEncoder)
	}

	e.Value = make([]byte, len(value))
	copy(e.Value, value)
	read += len(e.Value)

	return read, nil
}

func (e *Element) binaryDecoding(raw []byte) (int, error) {
	var value []byte
	var err error
	var contentLen int
	var read int

	contentLen = e.Length
	if e.Encoding == utils.EncodingChar {
		if len(raw) < read+contentLen {
			return 0, errors.New(utils.ErrBadElementData)
		}

		value, err = utils.UTF8ToWindows1252(raw[read : read+contentLen])
		if err != nil {
			return 0, err
		}
		read += contentLen
	} else if e.Encoding == utils.EncodingHex {
		if len(raw) < read+contentLen/4 {
			return 0, errors.New(utils.ErrBadElementData)
		}
		hexNumber, err := strconv.ParseUint(string(raw[read:read+contentLen/4]), 16, contentLen)
		if err != nil {
			return 0, err
		}

		binaryNumber := strconv.FormatUint(hexNumber, 2)
		formatStr := "%0" + strconv.Itoa(e.Length) + "s"
		binaryNumber = fmt.Sprintf(formatStr, binaryNumber)
		e.Value = make([]byte, len(binaryNumber))
		copy(e.Value, binaryNumber)
		e.extendBinaryData()
		read += contentLen / 4

		return read, nil
	} else {
		return 0, errors.New(utils.ErrInvalidEncoder)
	}

	e.Value = make([]byte, len(value))
	copy(e.Value, value)

	return read, nil
}

func (e *Element) lengthEncoding(value []byte) ([]byte, error) {
	var encode []byte
	lenStr := strconv.Itoa(e.Length)
	formatStr := "%0" + strconv.Itoa(len(lenStr)) + "d"
	contentLen := []byte(fmt.Sprintf(formatStr, len(value)))

	switch e.LengthEncoding {
	case utils.EncodingChar:
		encode = contentLen
	case utils.EncodingHex:
		lenStr = strconv.Itoa(len(fmt.Sprintf("%x", e.Length)))
		formatStr = "%0" + strconv.Itoa(len(lenStr)) + "x"
		encode = []byte(fmt.Sprintf(formatStr, len(value)))
	case utils.EncodingRBcd:
		// contentLen is number only
		encode, _ = utils.RBcd(contentLen)
	case utils.EncodingBcd:
		// contentLen is number only
		encode, _ = utils.Bcd(contentLen)
	default:
		return nil, errors.New(utils.ErrInvalidLengthEncoder)
	}

	return encode, nil
}

func (e *Element) lengthDecoding(raw []byte) (int, error) {
	var contentLen int
	var read int
	var err error

	if !e.Fixed {
		lenSize := len(strconv.Itoa(e.Length))
		bcdSize := lenSize / 2
		if lenSize%2 != 0 {
			bcdSize++
		}

		if len(raw) < lenSize {
			return 0, errors.New(utils.ErrBadElementData)
		}

		switch e.LengthEncoding {
		case utils.EncodingChar:
			contentLen, err = strconv.Atoi(string(raw[:lenSize]))
			if err != nil {
				return 0, errors.New(utils.ErrParseLengthFailed + ": " + string(raw[:lenSize]))
			}
			read = lenSize
		case utils.EncodingHex:
			lenStr := strconv.Itoa(len(fmt.Sprintf("%x", e.Length)))
			read = len(lenStr)
			n, err := strconv.ParseInt(string(raw[:read]), 16, 32)
			if err != nil {
				return 0, errors.New(utils.ErrParseLengthFailed + ": " + string(raw[:read]))
			}
			contentLen = int(n)
		case utils.EncodingRBcd:
			lenVal, err := utils.RBcdAscii(raw[:bcdSize], lenSize)
			if err != nil {
				return 0, err
			}
			// lenVal is decimal numbers only
			contentLen, _ = strconv.Atoi(string(lenVal))
			read = bcdSize
		case utils.EncodingBcd:
			lenVal, err := utils.BcdAscii(raw[:bcdSize], lenSize)
			if err != nil {
				return 0, err
			}
			// lenVal is decimal numbers only
			contentLen, _ = strconv.Atoi(string(lenVal))
			read = bcdSize
		default:
			return 0, errors.New(utils.ErrInvalidLengthEncoder)
		}
		e.DataLength = contentLen
	}

	return read, nil
}

func (e *Element) extendBinaryData() {
	cat := utils.AvailableTypeCategory[e.Type]
	if cat == utils.EncodingCatBinary && (len(e.Value) < e.Length) {
		newData := fmt.Sprintf("%-"+strconv.Itoa(e.Length)+"s", string(e.Value))
		newData = strings.ReplaceAll(newData, " ", "0")
		e.Value = make([]byte, e.Length)
		copy(e.Value, newData)
	}
}

func (e *Element) setType(_type *utils.ElementType) {
	e.Type = _type.Type
	e.Length = _type.Length
	e.Format = _type.Format
	e.Encoding = _type.Encoding
	e.Fixed = _type.Fixed
	e.LengthEncoding = _type.LengthEncoding
	e.extendBinaryData()
}

func (e *Element) numberWithPadding(buf []byte, isNumeric bool) []byte {
	var length = e.Length
	if !e.Fixed {
		return buf
	}
	if e.Encoding == utils.EncodingChar {
		return []byte(fmt.Sprintf("%0"+strconv.Itoa(length)+"s", buf))
	}

	bcdSize := length / 2
	if length%2 != 0 {
		bcdSize++
	}
	value := make([]byte, bcdSize)
	if e.Encoding == utils.EncodingBcd {
		for index := range value {
			if index < len(buf) {
				value[index] = buf[index]
			} else {
				value[index] = 0
			}
		}
	} else {
		pos := bcdSize - len(buf)
		for index := range value {
			if index >= pos {
				value[index] = buf[index-pos]
			} else {
				value[index] = 0
			}
		}
	}
	return value
}

func (e *Element) characterWithPadding(buf []byte) []byte {
	var length = e.Length
	if !e.Fixed {
		return buf
	}
	fmtStr := "%-" + strconv.Itoa(length) + "s"
	return []byte(fmt.Sprintf(fmtStr, buf))
}

func (e *Element) validateWithRegex() error {
	// Regex
	match := false
	if len(e.Value) == 0 {
		return nil
	}
	switch e.Type {
	case utils.ElementTypeAlphabetic:
		match = utils.RegexAlphabetic(string(e.Value))
	case utils.ElementTypeNumeric, utils.ElementTypeMti:
		match = utils.RegexNumeric(string(e.Value))
	case utils.ElementTypeSpecial:
		match = utils.RegexSpecial(string(e.Value))
	case utils.ElementTypeIndicate:
		match = utils.RegexIndicate(string(e.Value))
	case utils.ElementTypeBinary, utils.ElementTypeBitmap:
		match = utils.RegexBinary(string(e.Value))
	case utils.ElementTypeAlphaNumeric:
		match = utils.RegexAlphaNumeric(string(e.Value))
	case utils.ElementTypeAlphaSpecial:
		match = utils.RegexAlphaSpecial(string(e.Value))
	case utils.ElementTypeNumericSpecial:
		match = utils.RegexNumericSpecial(string(e.Value))
	case utils.ElementTypeAlphaNumericSpecial:
		match = utils.RegexAlphaNumericSpecial(string(e.Value))
	case utils.ElementTypeIndicateNumeric:
		match = utils.RegexIndicateNumeric(string(e.Value))
	}
	if !match {
		return errors.New(utils.ErrBadElementData)
	}
	return nil
}

func (e *Element) validateWithEncoding() error {
	// Checking available encoding
	available := utils.CheckAvailableEncoding(e.Type, e.Encoding)
	if !available {
		return errors.New(utils.ErrNonAvailableEncoding)
	}

	// Checking available encoding of length
	if !e.Fixed {
		available := utils.CheckAvailableEncoding(utils.ElementTypeNumberEncoding, e.LengthEncoding)
		if !available {
			return errors.New(utils.ErrNonAvailableEncoding)
		}
	}
	return nil
}

func (e *Element) validateWithFormat() error {
	if len(e.Format) == 0 {
		return nil
	}
	format := strings.ToUpper(e.Format)
	f, exist := utils.AvailableDateFormat[format]
	if exist {
		match := f(string(e.Value))
		if !match {
			return errors.New(utils.ErrBadElementData)
		}
	}
	return nil
}
