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
	assert.Equal(t, nil, err)

	err = json.Unmarshal(jsonStr, &jsonMessage)
	assert.Equal(t, nil, err)

	orgJsonBuf, err := json.MarshalIndent(&jsonMessage, "", "\t")
	assert.Equal(t, nil, err)

	orgXmlBuf, err := xml.MarshalIndent(&jsonMessage, "", "\t")
	assert.Equal(t, nil, err)

	xmlMessage, err := NewDataElements(&utils.ISO8583DataElementsVer1987)
	assert.Equal(t, nil, err)

	err = xml.Unmarshal(xmlStr, &xmlMessage)
	assert.Equal(t, nil, err)

	jsonBuf, err := json.MarshalIndent(&xmlMessage, "", "\t")
	assert.Equal(t, nil, err)

	xmlBuf, err := xml.MarshalIndent(&xmlMessage, "", "\t")
	assert.Equal(t, nil, err)

	assert.Equal(t, orgJsonBuf, jsonBuf)
	assert.Equal(t, orgXmlBuf, xmlBuf)
}
