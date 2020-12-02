// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/moov-io/iso8583/pkg/utils"
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

	rawMessage := &dataElements{}
	err = json.Unmarshal(jsonStr, rawMessage)
	assert.NotNil(t, err)

	err = xml.Unmarshal(xmlStr, rawMessage)
	assert.NotNil(t, err)

	rawMessage.spec = &utils.ISO8583DataElementsVer1987
	err = json.Unmarshal(jsonStr, rawMessage)
	assert.Nil(t, err)

	invalidJsonStr := []byte(`{
		"3": 123456,
		"11": 123456,
		"238": "abcdef"
	}`)
	err = json.Unmarshal(invalidJsonStr, rawMessage)
	assert.NotNil(t, err)

	invalidXmlStr := []byte(`<DataElements>
		<Element Number="3">123456</Element>
		<Element Number="11">123456</Element>
		<Element Number="238">abcdef</Element>
	</DataElements>`)
	err = xml.Unmarshal(invalidXmlStr, rawMessage)
	assert.NotNil(t, err)

	rawMessage.spec = &utils.Specification{
		Encoding: &utils.EncodingDefinition{
			MtiEnc:    utils.EncodingChar,
			BitmapEnc: utils.EncodingHex,
			BinaryEnc: utils.EncodingChar,
		},
		Elements: &utils.Attributes{
			3: {Describe: "n6", Description: "Processing code"},
		},
	}
	invalidJsonStr = []byte(`{
		"3": 123456
	}`)
	err = json.Unmarshal(invalidJsonStr, rawMessage)
	assert.NotNil(t, err)

	invalidJsonStr = []byte(`{
		"3": 123456,
	}`)
	err = json.Unmarshal(invalidJsonStr, rawMessage)
	assert.NotNil(t, err)

	invalidXmlStr = []byte(`<DataElements>
		<Element Number="3">123456</Element>
	</DataElements>`)
	err = xml.Unmarshal(invalidXmlStr, rawMessage)
	assert.NotNil(t, err)

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

	message, err := NewISO8583Message(&utils.ISO8583DataElementsVer1987)
	assert.Nil(t, err)

	err = json.Unmarshal(jsonStr, message)
	assert.Nil(t, err)

	_, err = json.MarshalIndent(message, "", "\t")
	assert.Nil(t, err)

	_, err = xml.MarshalIndent(message, "", "\t")
	assert.Nil(t, err)

	jsonStr = []byte(`
	{
		"mti": "0800",
		"bitmap": "10100000001000000000000000000000000001",
		"elements": {
			"1": "00000000000000000000000000000000000",
			"11": 123456,
			"3": 123456,
			"asdf": "abcdef"
		}
	}
	`)
	err = json.Unmarshal(jsonStr, message)
	assert.NotNil(t, err)
}

func TestElementStruct(t *testing.T) {
	element := &Element{}
	element.Value = []byte("123456")
	element.Length = 6
	element.Encoding = utils.EncodingChar
	element.LengthEncoding = utils.EncodingChar

	element.Type = utils.ElementTypeAlphabetic
	element.Encoding = utils.EncodingAscii
	err := element.Validate()
	assert.NotNil(t, err)

	_, err = element.Load(nil)
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeNumeric
	element.Encoding = utils.EncodingChar
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeMti
	err = element.Validate()
	assert.Nil(t, err)

	element.Type = utils.ElementTypeSpecial
	element.Encoding = utils.EncodingAscii
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeIndicate
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeBinary
	element.Encoding = utils.EncodingChar
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeBitmap
	err = element.Validate()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeAlphaNumeric
	element.Encoding = utils.EncodingAscii
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
	assert.Nil(t, err)

	// cat number, fixed
	element.DataLength = 6
	element.Fixed = true
	buf, err := element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte("123456"))

	element.Encoding = "unknown"
	element.Fixed = true
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingRBcd
	element.Type = "unknown"
	_, err = element.Bytes()
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeNumeric
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x12, 0x34, 0x56})

	element.Value = []byte("123")
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x0, 0x1, 0x23})
	element.Value = []byte("123456")

	element.Encoding = utils.EncodingBcd
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x12, 0x34, 0x56})

	element.Value = []byte("123")
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x12, 0x30, 0x0})
	element.Value = []byte("123456")

	element.Fixed = false
	element.Encoding = utils.EncodingChar
	element.LengthEncoding = utils.EncodingChar
	buf, err = element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte("6123456"))

	element.LengthEncoding = utils.EncodingAscii
	_, err = element.Bytes()
	assert.NotNil(t, err)

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

	element.LengthEncoding = utils.EncodingChar
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	buf = []byte("06abcdef")
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.LengthEncoding = "unknown"
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingRBcd
	element.Type = "unknown"
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.Type = utils.ElementTypeAlphabetic
	_, err = element.Load(buf)
	assert.Nil(t, err)

	element.LengthEncoding = utils.EncodingBcd
	_, err = element.Load(buf)
	assert.Nil(t, err)

	buf = []byte("abcdef")
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
	element.Fixed = false

	_, err = element.Load(nil)
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingAscii
	buf = []byte("123456")
	_, err = element.Load(buf)
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingChar
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

	element.Value = []byte("123456")
	element.Length = 6
	element.Type = utils.ElementTypeNumeric
	_, err = xml.Marshal(element)
	assert.Nil(t, err)

	element.Type = utils.ElementTypeAlphabetic
	_, err = xml.Marshal(element)
	assert.Nil(t, err)

	element = &Element{}
	element.Type = utils.ElementTypeAlphabetic
	err = xml.Unmarshal([]byte(`<Element>123456</Element>`), element)
	assert.Nil(t, err)
}

func TestElementStructForErrorCases(t *testing.T) {
	element := &Element{}
	element.Value = []byte("1234567")
	element.Length = 7
	element.Fixed = true
	element.Type = utils.ElementTypeNumeric
	element.Encoding = utils.EncodingBcd
	buf, err := element.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, []byte{0x12, 0x34, 0x56, 0x70})

	element.Length = 7
	element.Value = nil

	element.Fixed = false
	element.Type = utils.ElementTypeNumeric
	element.LengthEncoding = utils.EncodingHex
	element.Encoding = utils.EncodingBcd
	element.Type = utils.ElementTypeAlphabetic
	_, err = element.Load([]byte{0x7, 0x12, 0x34, 0x56, 0x70})
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingRBcd
	_, err = element.Load([]byte{0xA, 0x12, 0x34, 0x56, 0x70})
	assert.NotNil(t, err)

	element.LengthEncoding = utils.EncodingBcd
	_, err = element.Load([]byte{0xA, 0x12, 0x34, 0x56, 0x70})
	assert.NotNil(t, err)

	element.Fixed = true
	element.Encoding = utils.EncodingEbcdic
	_, err = element.Load([]byte{0xE3, 0x88, 0x89, 0xA2, 0x40, 0xA2})
	assert.NotNil(t, err)

	element.Fixed = false
	element.LengthEncoding = utils.EncodingHex
	element.DataLength = 6
	element.Value = []byte("123456")
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.LengthEncoding = utils.EncodingHex
	element.DataLength = 0
	element.Value = nil
	_, err = element.Bytes()
	assert.Nil(t, err)

	element.Length = 9
	element.Fixed = true
	element.Type = utils.ElementTypeBinary
	element.Encoding = utils.EncodingHex
	_, err = element.Load([]byte("A"))
	assert.NotNil(t, err)

	element.Length = 9
	element.Encoding = utils.EncodingChar
	element.Type = utils.ElementTypeNumeric
	_, err = element.Load([]byte("A"))
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingRBcd
	_, err = element.Load([]byte("A"))
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingBcd
	_, err = element.Load([]byte("A"))
	assert.NotNil(t, err)

	element.Length = 13
	element.Encoding = utils.EncodingRBcd
	_, err = element.Load([]byte("unacceptable"))
	assert.NotNil(t, err)

	element.Encoding = utils.EncodingBcd
	_, err = element.Load([]byte("unacceptable"))
	assert.NotNil(t, err)

	element.Length = 9
	element.Fixed = true
	element.Type = utils.ElementTypeBinary
	element.Encoding = utils.EncodingHex
	element.Value = []byte("0101")
	_, err = element.Bytes()
	assert.NotNil(t, err)
}

func TestIso8583MessageBytes(t *testing.T) {
	byteData := []byte(`0800a0200000040000000000000000000000123456123456abcdef`)

	message, err := NewISO8583Message(&utils.ISO8583DataElementsVer1987)
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

	_, err = message.Bytes()
	assert.Nil(t, err)

	_element := mapRet[3]
	copy(_element.Value, "abcd")
	mapRet[3] = _element
	err = message.Validate()
	assert.NotNil(t, err)

	ret = message.GetBitmap()
	copy(ret.Value, "abcd0001")
	err = message.Validate()
	assert.NotNil(t, err)

	ret = message.GetMti()
	copy(ret.Value, "ABCD")
	err = message.Validate()
	assert.NotNil(t, err)

	//byteData = []byte(`0800a020000004000000000000000000000000000`)
	byteData = []byte(`ABC`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	byteData = []byte(`0800PPP0000004000000000000000000000000000`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	byteData = []byte(`08000000000000000000000000000000000000000EOF`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	byteData = []byte(`08000000000000000000`)
	_, err = message.Load(byteData)
	assert.Nil(t, err)

	_, err = NewISO8583Message(nil)
	assert.NotNil(t, err)

	ret = message.GetBitmap()
	copy(ret.Value, "F010000000100000000000000000000000000000000000000000000000000000")
	err = message.Validate()
	assert.NotNil(t, err)

	_, err = message.Bytes()
	assert.NotNil(t, err)

	ret = message.GetMti()
	copy(ret.Value, "FFFF")
	ret.Encoding = utils.EncodingRBcd
	_, err = message.Bytes()
	assert.NotNil(t, err)

	mapRet = message.GetElements()
	modified := mapRet[11]
	copy(modified.Value, "FFFF")
	modified.Encoding = utils.EncodingRBcd
	message = &isoMessage{elements: &dataElements{
		elements: mapRet,
		spec:     &utils.ISO8583DataElementsVer1987,
	}}
	_, err = message.Bytes()
	assert.NotNil(t, err)

	message = &isoMessage{elements: nil}
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	mapRet = message.GetElements()
	assert.Equal(t, len(mapRet), 0)

	xmlStr := []byte(`<DataElements>
		<Element Number="3">123456</Element>
		<Element Number="11">123456</Element>
		<Element Number="38">abcdef</Element>
	</DataElements>`)
	err = xml.Unmarshal(xmlStr, message)
	assert.Nil(t, err)

	mapRet = message.GetElements()
	assert.Equal(t, len(mapRet), 0)

	_spec := &utils.Specification{
		Encoding: &utils.EncodingDefinition{
			MtiEnc:    utils.EncodingChar,
			BitmapEnc: utils.EncodingHex,
			BinaryEnc: utils.EncodingChar,
		},
		Elements: &utils.Attributes{
			1: {Describe: "b 64", Description: "Second Bitmap"},
			2: {Describe: "b 64", Description: "Second Bitmap"},
		},
	}
	message, err = NewISO8583Message(_spec)
	assert.Nil(t, err)
	byteData = []byte(
		`0800c000000000000000` +
			`0000000000000000000000000000000000000000000000000000000000000000` +
			`0000000000000000000000000000000000000000000000000000000000000000`)
	_, err = message.Load(byteData)
	assert.Nil(t, err)

	byteData = []byte(`0800a020000004000000`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	byteData = []byte(`0800a0200000040000000000000000000000123456123456abcdef`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	byteData = []byte(
		`0800E000000000000000` +
			`0000000000000000000000000000000000000000000000000000000000000000` +
			`0000000000000000000000000000000000000000000000000000000000000000` +
			`0000000000000000000000000000000000000000000000000000000000000000`)
	_, err = message.Load(byteData)
	assert.NotNil(t, err)

	_spec = &utils.Specification{
		Encoding: &utils.EncodingDefinition{
			MtiEnc:    utils.EncodingChar,
			BitmapEnc: utils.EncodingHex,
			BinaryEnc: utils.EncodingChar,
		},
		Elements: &utils.Attributes{
			1: {Describe: "b 64", Description: "Second Bitmap"},
			2: {Describe: "n64", Description: "Number"},
		},
	}
	message, err = NewISO8583Message(_spec)
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.NotNil(t, err)
}

func TestElementsStruct(t *testing.T) {
	_, err := NewDataElements(nil)
	assert.NotNil(t, err)

	message, err := NewDataElements(&utils.ISO8583DataElementsVer1987)
	assert.Nil(t, err)

	message.elements[1] = &Element{
		Type:           utils.ElementTypeNumeric,
		Length:         4,
		Encoding:       "",
		Fixed:          true,
		LengthEncoding: "",
		Value:          []byte("abcd"),
	}
	err = message.Validate()
	assert.NotNil(t, err)

	_, err = message.Bytes()
	assert.NotNil(t, err)

	_, err = json.Marshal(message)
	assert.NotNil(t, err)

	message.elements[1] = &Element{
		Type:           utils.ElementTypeAlphabetic,
		Length:         4,
		Encoding:       utils.EncodingAscii,
		Fixed:          true,
		LengthEncoding: "",
		Value:          []byte("abcd"),
	}
	_, err = message.Bytes()
	assert.Nil(t, err)

	_, err = message.Load(nil)
	assert.Nil(t, err)
}

func TestISO8583MessageWithValidSamples(t *testing.T) {
	samples := []string{
		"financial_transaction_message.dat",
		"financial_transaction_message_response.dat",
		"iso_reversal_message.dat",
		"iso_reversal_message_response.dat",
		"iso_reversal_repeat_message.dat",
		"iso_reversal_repeat_message_response.dat",
		"network_management_message.dat",
		"network_management_message_response.dat",
		"network_management_message_with_track.dat",
	}

	for _, sample := range samples {
		message, err := NewISO8583Message(&utils.ISO8583DataElementsVer1987)
		assert.Nil(t, err)

		byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", sample))
		assert.Nil(t, err)

		_, err = message.Load(byteData)
		assert.Nil(t, err)

		err = message.Validate()
		assert.Nil(t, err)

		buf, err := message.Bytes()
		assert.Nil(t, err)
		assert.Equal(t, buf, byteData)
	}

	samples = []string{
		"iso_reversal_message_error_date.dat",
		"network_management_message_with_error_track.dat",
	}
	for _, sample := range samples {
		message, err := NewISO8583Message(&utils.ISO8583DataElementsVer1987)
		assert.Nil(t, err)

		byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", sample))
		assert.Nil(t, err)

		_, err = message.Load(byteData)
		assert.Nil(t, err)

		err = message.Validate()
		assert.NotNil(t, err)

		buf, err := message.Bytes()
		assert.Nil(t, err)
		assert.Equal(t, buf, byteData)
	}
}

func TestISO8583MessageWithJson(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	_, err = NewISO8583MessageWithJson(jsonData, nil)
	assert.Nil(t, err)

	jsonData = []byte(`{
	"1": {
		"Describe": "b 64",
		"Description": "Second Bitmap"
	},}`)
	_, err = NewISO8583MessageWithJson(jsonData, nil)
	assert.NotNil(t, err)
}

func TestISO8583MessageWithHexLengthEncoding(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingHex,
		NumberEnc:    utils.EncodingChar,
		CharacterEnc: utils.EncodingAscii,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_hex_length.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}

func TestISO8583MessageWithBcdLengthEncoding(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingBcd,
		NumberEnc:    utils.EncodingChar,
		CharacterEnc: utils.EncodingAscii,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_hex_bcd.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}

func TestISO8583MessageWithRBcdLengthEncoding(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingRBcd,
		NumberEnc:    utils.EncodingChar,
		CharacterEnc: utils.EncodingAscii,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_hex_rbcd.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}

func TestISO8583MessageForIndicateNumeric(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingRBcd,
		NumberEnc:    utils.EncodingChar,
		CharacterEnc: utils.EncodingAscii,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_indicate_numeric.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)

	byteData, err = ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_error_indicate_numeric.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.NotNil(t, err)

	buf, err = message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}

func TestISO8583MessageWithBcdNumberEncoding(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingRBcd,
		NumberEnc:    utils.EncodingBcd,
		CharacterEnc: utils.EncodingAscii,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_number_bcd_encoding.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}

func TestISO8583MessageWithRBcdNumberEncoding(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingRBcd,
		NumberEnc:    utils.EncodingRBcd,
		CharacterEnc: utils.EncodingAscii,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_number_bcd_encoding.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}

func TestISO8583MessageWithCharacterEbcdicEncoding(t *testing.T) {
	jsonData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "attributes_data_elements.dat"))
	assert.Nil(t, err)

	encoding := &utils.EncodingDefinition{
		MtiEnc:       utils.EncodingChar,
		BitmapEnc:    utils.EncodingHex,
		LengthEnc:    utils.EncodingChar,
		NumberEnc:    utils.EncodingChar,
		CharacterEnc: utils.EncodingEbcdic,
		BinaryEnc:    utils.EncodingHex,
		TrackEnc:     utils.EncodingEbcdic,
	}
	message, err := NewISO8583MessageWithJson(jsonData, encoding)
	assert.Nil(t, err)

	byteData, err := ioutil.ReadFile(filepath.Join("..", "..", "test", "testdata", "message_with_character_ebcdic_encoding.dat"))
	assert.Nil(t, err)

	_, err = message.Load(byteData)
	assert.Nil(t, err)

	err = message.Validate()
	assert.Nil(t, err)

	buf, err := message.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, buf, byteData)
}
