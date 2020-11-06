// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package pkg

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Attribute struct {
	Describe    string // [attribute(b 64, b-64, b..64)]; [format(MMDD, hhmmss)]
	Description string
}

type EncodingDefinition struct {
	MtiEnc       string `json:"mti_enc"`
	BitmapEnc    string `json:"bmp_enc"`
	LengthEnc    string `json:"len_enc"`
	NumberEnc    string `json:"num_enc"`
	CharacterEnc string `json:"chr_enc"`
	BinaryEnc    string `json:"bin_enc"`
	TrackEnc     string `json:"trk_enc"`
}

var (
	digitIndicates = []string{
		"...",
		"..",
		".",
		"-",
		" ",
	}
	formatIndicate   = ";"
	variableIndicate = "."
)

func (s Attribute) ElementType() (*CommonType, error) {
	attribute := s.Describe
	for _, indicate := range digitIndicates {
		if strings.Contains(attribute, indicate) {
			var format string
			if strings.Contains(attribute, formatIndicate) {
				splits := strings.Split(attribute, formatIndicate)
				attribute = splits[0]
				format = strings.TrimSpace(splits[len(splits)-1])
			}

			isFixed := !strings.Contains(indicate, variableIndicate)
			splits := strings.Split(attribute, indicate)
			_size, err := strconv.Atoi(strings.TrimSpace(splits[len(splits)-1]))
			if err != nil ||
				(!isFixed && _size > int(math.Pow(10, float64(len(indicate))))) {
				return nil, errors.New("invalid element length")
			}

			return &CommonType{
				Type:   strings.TrimSpace(splits[0]),
				Length: _size,
				Fixed:  isFixed,
				Format: format,
			}, nil
		}
	}

	return nil, errors.New("invalid element type")
}

type Attributes map[int]Attribute
type Specification struct {
	Elements *Attributes         `json:"elements,omitempty"`
	Encoding *EncodingDefinition `json:"encoding,omitempty"`
}

func (s Attributes) Keys() (keys []int) {
	for k, _ := range s {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return
}

func (s Attributes) Get(number int) (*Attribute, error) {
	spec, existed := s[number]
	if !existed {
		return nil, errors.New("don't exist specification")
	}
	return &spec, nil
}
