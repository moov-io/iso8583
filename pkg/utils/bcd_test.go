// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
