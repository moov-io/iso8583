// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"github.com/yerden/go-util/bcd"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var standard = &bcd.BCD{
	Map: map[byte]byte{
		'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
		'4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
		'8': 0x8, '9': 0x9,
	},
	SwapNibbles: false,
	Filler:      0x0}

// Bcd character string to bcd string
func Bcd(src []byte) ([]byte, error) {
	enc := bcd.NewEncoder(standard)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// RBcd character string to right aligned bcd string
func RBcd(src []byte) ([]byte, error) {
	if len(src)%2 != 0 {
		src = append([]byte("0"), src...)
	}
	return Bcd(src)
}

// BcdAscii bcd string to ascii string
func BcdAscii(src []byte, length int) ([]byte, error) {
	dec := bcd.NewDecoder(standard)
	dst := make([]byte, bcd.DecodedLen(len(src)))
	n, err := dec.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	if n > length {
		n = length
	}
	return dst[:n], err
}

// RBcdAscii right aligned bcd string to ascii string
func RBcdAscii(src []byte, length int) ([]byte, error) {
	dec := bcd.NewDecoder(standard)
	dst := make([]byte, bcd.DecodedLen(len(src)))
	n, err := dec.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	start := n - length
	if start < 0 {
		start = 0
	}
	return dst[start:n], err
}

// UTF8ToWindows1252 converts text encoded in UTF-8 to Windows-1252 or CP-1252 encoding
func UTF8ToWindows1252(input []byte) ([]byte, error) {
	if utf8.Valid(input) {
		reader := bytes.NewReader(input)
		res, err := transformEncoding(reader, charmap.Windows1252.NewEncoder())
		return res, err
	}
	return input, nil
}

// HexToBitmapArray converts a hex string to a bit array
func BitmapToIndexArray(bitmap string, base int) []int {
	bitArrayStrings := strings.Split(bitmap, "")
	bitArray := make([]int, 0)
	for index := 0; index < len(bitArrayStrings); index++ {
		if bitArrayStrings[index] == "1" {
			bitArray = append(bitArray, index+base+1)
		}
	}
	return bitArray
}

// IsSecondBitmap return existence of sub bitmap
func IsSecondBitmap(bitmap string) bool {
	indexes := BitmapToIndexArray(bitmap, 0)
	if len(indexes) > 0 && indexes[0] == 1 {
		return true
	}
	return false
}

// IsThirdBitmap return existence of sub bitmap
func IsThirdBitmap(bitmap string) bool {
	indexes := BitmapToIndexArray(bitmap, 0)
	if len(indexes) > 1 && indexes[0] == 1 && indexes[1] == 2 {
		return true
	}
	return false
}

// Get message format
func MessageFormat(buf []byte) string {
	if isValidJSON(buf) {
		return MessageFormatJson
	} else if isValidXML(buf) {
		return MessageFormatXml
	}
	return MessageFormatIso8583
}

func isValidXML(buf []byte) bool {
	decoder := xml.NewDecoder(bytes.NewBuffer(buf))
	err := xml.Unmarshal(buf, new(interface{}))
	if err != nil {
		return false
	}
	for {
		err = decoder.Decode(new(interface{}))
		if err != nil {
			break
		}
	}
	return err == io.EOF
}

func isValidJSON(buf []byte) bool {
	var dummy map[string]interface{}
	if err := json.Unmarshal(buf, &dummy); err != nil {
		return false
	}
	return true
}

func transformEncoding(reader io.Reader, trans transform.Transformer) ([]byte, error) {
	transReader := transform.NewReader(reader, trans)
	ret, err := ioutil.ReadAll(transReader)
	return ret, err
}
