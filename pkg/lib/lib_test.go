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
		"bitmap": "823A000000000000040000000000000004200906139000010906130420042000",
		"message": {
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
