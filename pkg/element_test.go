package pkg

import (
	"encoding/json"
	"encoding/xml"
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

	jsonMessage := NewDataElements(ISO8583DataElementsVer1987)
	err := json.Unmarshal(jsonStr, &jsonMessage)
	assert.Equal(t, nil, err)

	orgJsonBuf, err := json.MarshalIndent(&jsonMessage, "", "\t")
	assert.Equal(t, nil, err)

	_, err = xml.MarshalIndent(&jsonMessage, "", "\t")
	assert.Equal(t, nil, err)

	xmlMessage := NewDataElements(ISO8583DataElementsVer1987)
	err = xml.Unmarshal(xmlStr, &xmlMessage)
	assert.Equal(t, nil, err)

	jsonBuf, err := json.MarshalIndent(&xmlMessage, "", "\t")
	assert.Equal(t, nil, err)

	_, err = xml.MarshalIndent(&xmlMessage, "", "\t")
	assert.Equal(t, nil, err)

	assert.Equal(t, orgJsonBuf, jsonBuf)
}
