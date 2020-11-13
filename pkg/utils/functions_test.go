// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAvailableEncoding(t *testing.T) {
	b := CheckAvailableEncoding(ElementTypeMti, EncodingBcd)
	assert.Equal(t, b, true)

	b = CheckAvailableEncoding(ElementTypeMti, EncodingAscii)
	assert.Equal(t, b, false)
}

func TestBcdEncode(t *testing.T) {
	b := []byte("954")
	r, err := RBcd(b)
	assert.Nil(t, err)
	assert.Equal(t, "0954", fmt.Sprintf("%X", r))

	r, err = Bcd(b)
	assert.Nil(t, err)
	assert.Equal(t, "9540", fmt.Sprintf("%X", r))

	b = []byte("31")
	r, err = Bcd(b)
	assert.Nil(t, err)
	assert.Equal(t, "31", fmt.Sprintf("%X", r))

	r, err = RBcd(b)
	assert.Nil(t, err)
	assert.Equal(t, "31", fmt.Sprintf("%X", r))

}

func TestBcdDecode(t *testing.T) {
	_, err := BcdAscii([]byte("\x12\xa3\x4f"), 6)
	assert.NotNil(t, err)

	r, err := BcdAscii([]byte("\x12\x34\x56"), 6)
	assert.Nil(t, err)
	assert.Equal(t, []byte("123456"), r)

	r, err = BcdAscii([]byte("\x12\x04\x50"), 5)
	assert.Nil(t, err)
	assert.Equal(t, []byte("12045"), r)

	r, err = RBcdAscii([]byte("\x01\x23\x45"), 5)
	assert.Nil(t, err)
	assert.Equal(t, []byte("12345"), r)
}

func TestAttributeParse(t *testing.T) {
	attributes := ISO8583DataElementsVer1987.Elements
	for _, key := range attributes.Keys() {
		attr, err := attributes.Get(key)
		assert.Nil(t, err)
		_, err = attr.Parse()
		assert.Nil(t, err)
	}

	_, err := attributes.Get(-1)
	assert.NotNil(t, err)
	attribute := Attribute{
		Description: "Function code (ISO 8583:1993), or network international identifier (NII)",
		Describe:    "n 3n",
	}
	_, err = attribute.Parse()
	assert.NotNil(t, err)
}

func TestUTF8ToWindows1252(t *testing.T) {
	test := "test"
	convert, err := UTF8ToWindows1252([]byte(test))
	assert.Nil(t, err)
	assert.Equal(t, test, string(convert))
}

func TestBitmapToIndexArray(t *testing.T) {
	bitmap := "1000000000100001"
	indexes := BitmapToIndexArray(bitmap, 0)
	assert.Equal(t, indexes, []int{1, 11, 16})
	indexes = BitmapToIndexArray(bitmap, 16)
	assert.Equal(t, indexes, []int{17, 27, 32})
}

func TestIsExistedBitmap(t *testing.T) {
	existed := IsSecondBitmap("1000000000100001")
	assert.Equal(t, existed, true)
	existed = IsSecondBitmap("0000000000100001")
	assert.Equal(t, existed, false)
	existed = IsThirdBitmap("1100000000100001")
	assert.Equal(t, existed, true)
	existed = IsThirdBitmap("1000000000100001")
	assert.Equal(t, existed, false)
}

func TestElementType(t *testing.T) {
	element := ElementType{Type: ElementTypeNumeric}
	err := element.Validate()
	assert.Nil(t, err)

	encoding := ISO8583DataElementsVer1987.Encoding
	element.SetEncoding(encoding)
	assert.Equal(t, element.Encoding, encoding.NumberEnc)

	element.Type = ElementTypeMti
	element.SetEncoding(encoding)
	assert.Equal(t, element.Encoding, encoding.MtiEnc)

	element.Type = ElementTypeBitmap
	element.SetEncoding(encoding)
	assert.Equal(t, element.Encoding, encoding.BitmapEnc)

	element.Type = ElementTypeBinary
	element.SetEncoding(encoding)
	assert.Equal(t, element.Encoding, encoding.BinaryEnc)

	element.Type = ElementTypeMagnetic
	element.SetEncoding(encoding)
	assert.Equal(t, element.Encoding, encoding.TrackEnc)

	element.Type = ElementTypeAlphabetic
	element.SetEncoding(encoding)
	assert.Equal(t, element.Encoding, encoding.CharacterEnc)
}
