// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
)

var (
	digitIndicates = []string{
		".....", // for variable data element (max 99999)
		"....",  // for variable data element (max 9999)
		"...",   // for variable data element (max 999)
		"..",    // for variable data element (max 99)
		".",     // for variable data element (max 9)
		"-",     // for fixed data element
		" ",     // for fixed data element
	}
	formatIndicate   = ";"
	variableIndicate = "."
)

type Attribute struct {
	Describe    string // [attribute(b 64, b-64, b..64)]; [format(MMDD, hhmmss)]
	Description string
}

// Parse return ElementType from attribute string
func (s Attribute) Parse() (*ElementType, error) {
	attribute := s.Describe
	for _, indicate := range digitIndicates {
		if strings.Contains(attribute, indicate) {
			var format string
			if strings.Contains(attribute, formatIndicate) {
				splits := strings.Split(attribute, formatIndicate)
				attribute = splits[0]
				if len(splits) > 1 {
					format = strings.TrimSpace(splits[len(splits)-1])
				}
			}
			splits := strings.Split(attribute, indicate)
			if len(splits) > 1 {
				isFixed := !strings.Contains(indicate, variableIndicate)
				_size, err := strconv.Atoi(strings.TrimSpace(splits[len(splits)-1]))
				if err != nil || (!isFixed && _size > int(math.Pow(10, float64(len(indicate))))) {
					return nil, errors.New(ErrInvalidElementLength)
				}
				return &ElementType{
					Type:   strings.TrimSpace(splits[0]),
					Length: _size,
					Fixed:  isFixed,
					Format: format,
				}, nil
			}
		}
	}

	return nil, errors.New(ErrInvalidElementType)
}

type Specification struct {
	Elements     *Attributes         `json:"elements,omitempty"`
	Encoding     *EncodingDefinition `json:"encoding,omitempty"`
	MessageTypes *MessageTypes       `json:"message_types,omitempty"`
}

type MessageType struct {
	MandatoryHexMask string `json:"mandatory_hex_mask,omitempty"`
	OptionalHexMask  string `json:"optional_hex_mask,omitempty"`
}

type Attributes map[int]Attribute
type MessageTypes map[string]MessageType

func (s Attributes) Keys() (keys []int) {
	for k := range s {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return
}

func (s Attributes) Get(number int) (*Attribute, error) {
	spec, existed := s[number]
	if !existed {
		return nil, errors.New(ErrNonExistSpecification)
	}
	return &spec, nil
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

// general element type for all of the data representation attributes
type ElementType struct {
	Type           string
	Length         int
	Format         string
	Encoding       string
	Fixed          bool
	LengthEncoding string
}

func (t *ElementType) Validate() error {
	return nil
}

// SetEncoding will set encoders
func (t *ElementType) SetEncoding(encoding *EncodingDefinition) {
	t.LengthEncoding = encoding.LengthEnc
	switch t.Type {
	case ElementTypeNumeric:
		t.Encoding = encoding.NumberEnc
	case ElementTypeMti:
		t.Encoding = encoding.MtiEnc
	case ElementTypeBitmap:
		t.Encoding = encoding.BitmapEnc
	case ElementTypeBinary:
		t.Encoding = encoding.BinaryEnc
	case ElementTypeMagnetic:
		t.Encoding = encoding.TrackEnc
	default:
		t.Encoding = encoding.CharacterEnc
	}
}
