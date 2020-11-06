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

type Specification struct {
	Attributes  string
	Description string
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

func (s Specification) ElementType() (*CommonType, error) {
	attributes := s.Attributes
	for _, indicate := range digitIndicates {
		if strings.Contains(attributes, indicate) {
			var format string
			if strings.Contains(attributes, formatIndicate) {
				splits := strings.Split(attributes, formatIndicate)
				attributes = splits[0]
				format = strings.TrimSpace(splits[len(splits)-1])
			}

			isFixed := !strings.Contains(indicate, variableIndicate)
			splits := strings.Split(attributes, indicate)
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

type Specifications map[int]Specification

func (s Specifications) Keys() (keys []int) {
	for k, _ := range s {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return
}

func (s Specifications) Get(number int) (*Specification, error) {
	spec, existed := s[number]
	if !existed {
		return nil, errors.New("don't exist specification")
	}
	return &spec, nil
}
