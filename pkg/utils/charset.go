// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"bytes"
	"io"
	"io/ioutil"
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
