// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"github.com/yerden/go-util/bcd"
)

var standard = &bcd.BCD{
	Map: map[byte]byte{
		'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
		'4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
		'8': 0x8, '9': 0x9,
	},
	SwapNibbles: false,
	Filler:      0x0}

func Bcd(src []byte) ([]byte, error) {
	enc := bcd.NewEncoder(standard)
	dst := make([]byte, bcd.EncodedLen(len(src)))
	n, err := enc.Encode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

func RBcd(src []byte) ([]byte, error) {
	if len(src)%2 != 0 {
		src = append([]byte("0"), src...)
	}
	return Bcd(src)
}

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
