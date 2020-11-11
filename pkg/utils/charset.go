// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func transformEnconding(reader io.Reader, trans transform.Transformer) ([]byte, error) {
	transReader := transform.NewReader(reader, trans)
	ret, err := ioutil.ReadAll(transReader)
	return ret, err
}

// UTF8ToWindows1252 converts text encoded in UTF-8 to Windows-1252 or CP-1252 encoding
func UTF8ToWindows1252(input []byte) ([]byte, error) {
	if utf8.Valid(input) {
		reader := bytes.NewReader(input)
		res, err := transformEnconding(reader, charmap.Windows1252.NewEncoder())
		return res, err
	}
	return input, nil
}

// BitMapArrayToHex converts a iso8583 bit array into a hex string
func BitMapArrayToHex(arr []int64) (string, error) {
	length := len(arr)
	m := make(map[float64]string)

	m[0] = "0"
	m[1] = "1"
	m[2] = "2"
	m[3] = "3"
	m[4] = "4"
	m[5] = "5"
	m[6] = "6"
	m[7] = "7"
	m[8] = "8"
	m[9] = "9"
	m[10] = "a"
	m[11] = "b"
	m[12] = "c"
	m[13] = "d"
	m[14] = "e"
	m[15] = "f"

	if (length % 4) != 0 {
		return "", errors.New(ErrInvalidBitmapArray)
	}

	if ((length / 4) % 2) != 0 {
		return "", errors.New(ErrInvalidBitmapArray)
	}
	var hexString string
	var buf float64
	var exp float64 = 3

	for index := 0; index < length; index++ {
		bit := arr[index] // get the bit at this index
		if bit == 0 {
			buf = buf + 0
			exp = exp - 1
		} else {
			buf = buf + math.Pow(2, exp)
			exp = exp - 1
		}

		// if exp is less than 0, it means we need to reset things
		if exp < 0 {
			exp = 3
			hexString = hexString + (m[buf])
			buf = 0
		}
	}

	return hexString, nil
}

// HexToBitmapArray converts a hex string to a bit array
func HexToBitmapArray(hexString string) ([]int64, error) {
	var bitString string
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	for index := 0; index < len(decoded); index++ {
		bitString = bitString + fmt.Sprintf("%8b", decoded[index])
	}
	bitArrayStrings := strings.Split(bitString, "")
	bitArray := make([]int64, len(bitArrayStrings))
	for index := 0; index < len(bitArrayStrings); index++ {
		bitArray[index], _ = strconv.ParseInt(bitArrayStrings[index], 10, 10)
	}
	return bitArray, nil
}
