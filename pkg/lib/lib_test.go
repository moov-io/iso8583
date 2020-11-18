// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"encoding/json"
	"encoding/xml"
	"github.com/moov-io/iso8583/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElementJsonXmlConvert(t *testing.T) {
	jsonStr := []byte(`{
		"3": 123456,
		"11": 123456,
		"38": "abcdef"
	}`)
	xmlStr := []byte(`<DataElements>
		<Element Number="3">123456</Element>
		<Element Number="11">123456</Element>
		<Element Number="38">abcdef</Element>
	</DataElements>`)

	jsonMessage, err := NewDataElements(&utils.ISO8583DataElementsVer1987)
	assert.Nil(t, err)

	err = json.Unmarshal(jsonStr, &jsonMessage)
	assert.Nil(t, err)

	orgJsonBuf, err := json.MarshalIndent(&jsonMessage, "", "\t")
	assert.Nil(t, err)

	orgXmlBuf, err := xml.MarshalIndent(&jsonMessage, "", "\t")
	assert.Nil(t, err)

	xmlMessage, err := NewDataElements(&utils.ISO8583DataElementsVer1987)
	assert.Nil(t, err)

	err = xml.Unmarshal(xmlStr, &xmlMessage)
	assert.Nil(t, err)

	jsonBuf, err := json.MarshalIndent(&xmlMessage, "", "\t")
	assert.Nil(t, err)

	xmlBuf, err := xml.MarshalIndent(&xmlMessage, "", "\t")
	assert.Nil(t, err)

	assert.Equal(t, orgJsonBuf, jsonBuf)
	assert.Equal(t, orgXmlBuf, xmlBuf)
}

func TestIso8583Message(t *testing.T) {
	jsonStr := []byte(`
	{
		"mti": "0800",
		"bitmap": "10100000001000000000000000000000000001",
		"elements": {
			"1": "00000000000000000000000000000000000",
			"11": 123456,
			"3": 123456,
			"38": "abcdef"
		}
	}
	`)

	message, err := NewMessage(&utils.ISO8583DataElementsVer1987)
	assert.Nil(t, err)

	err = json.Unmarshal(jsonStr, message)
	assert.Nil(t, err)

	_, err = json.MarshalIndent(message, "", "\t")
	assert.Nil(t, err)

	_, err = xml.MarshalIndent(message, "", "\t")
	assert.Nil(t, err)
}

func TestElementStruct(t *testing.T) {
	element := &Element{}
	element.Value = []byte("123456")
	element.Length = 6

	element.Type = utils.ElementTypeAlphabetic
	err := element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeNumeric
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeMti
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeSpecial
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeIndicate
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeBinary
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeBitmap
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeAlphaNumeric
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeAlphaSpecial
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeNumericSpecial
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeAlphaNumericSpecial
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeIndicateNumeric
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeMagnetic
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeNumeric
	number := element.String()
	assert.Equal(t, number, "123456")

	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingChar
	element.DataLength = 10
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.DataLength = 6
	element.Fixed = true
	buf, err := element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte("123456"))

	element.Encoding = utils.EncodingRBcd
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x12, 0x34, 0x56})

	buf, err = element.Bytes()
	element.Encoding = utils.EncodingBcd
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x12, 0x34, 0x56})

	element.Fixed = false
	element.Encoding = utils.EncodingChar
	element.LengthEncoding = utils.EncodingChar
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingAscii
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte("6123456"))

	element.LengthEncoding = utils.EncodingRBcd
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x06, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36})

	element.LengthEncoding = utils.EncodingBcd
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x60, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36})

	element.Value = []byte("abcdef")
	element.Type = utils.ElementTypeAlphaNumeric
	element.Encoding = utils.EncodingChar
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.Encoding = utils.EncodingAscii
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.Encoding = utils.EncodingEbcdic
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.Encoding = utils.EncodingHex
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.Value = []byte("1001001")
	element.Type = utils.ElementTypeBinary
	element.Encoding = utils.EncodingChar
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.Value = []byte("100100")
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.Encoding = utils.EncodingHex
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.Encoding = utils.EncodingAscii
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeAlphabetic
	element.LengthEncoding = utils.EncodingAscii
	buf = []byte("abcdef")
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingRBcd
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.LengthEncoding = utils.EncodingBcd
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingHex
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.Fixed = true
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Encoding = utils.EncodingEbcdic
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Encoding = utils.EncodingBcd
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeNumeric
	element.LengthEncoding = utils.EncodingAscii
	buf = []byte("123456")
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.LengthEncoding = utils.EncodingRBcd
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.LengthEncoding = utils.EncodingBcd
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.LengthEncoding = utils.EncodingHex
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Fixed = true
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Encoding = utils.EncodingEbcdic
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingChar
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Encoding = utils.EncodingRBcd
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Encoding = utils.EncodingBcd
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Type = utils.ElementTypeBinary
	element.Encoding = utils.EncodingAscii
	buf = []byte("100100")
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingChar
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.Encoding = utils.EncodingHex
	_, err = element.Load(buf)
	assert.Nil(t, err)
}

func TestIso8583MessageBytes(t *testing.T) {
	byteData := []byte(`0800a0200000040000000000000000000000000000000000000000000000000000000000000000000000123456123456abcdef`)

	message, err := NewMessage(&utils.ISO8583DataElementsVer1987)
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	ret := message.GetBitmap()
	assert.NotNil(t, ret)

	ret = message.GetMti()
	assert.NotNil(t, ret)

	mapRet := message.GetElements()
	assert.Equal(t, len(mapRet), 4)

	byteData = []byte(`0800a020000004000000000000000000000000000`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	_, err = NewMessage(nil)
	assert.NotNil(t, err)

	message = &isoMessage{elements: nil}
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	mapRet = message.GetElements()
	assert.Equal(t, len(mapRet), 0)
}
