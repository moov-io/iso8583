// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"errors"
	"github.com/moov-io/iso8583/pkg/utils"
)

// message is structure for ISO 8583 message encode and decode
type Message struct {
	Mti          *Element             `xml:"mti,omitempty" json:"mti,omitempty" yaml:"mti,omitempty"`
	Bitmap       *Element             `xml:"bitmap,omitempty" json:"bitmap,omitempty" yaml:"bitmap,omitempty"`
	Elements     *DataElements        `xml:"elements,omitempty" json:"message,elements" yaml:"elements,omitempty"`
	Spec         *utils.Specification `xml:"-" json:"-"`
	SecondBitmap bool                 `xml:"-" json:"-"`
	ThirdBitmap  bool                 `xml:"-" json:"-"`
}

// create data elements of message with specification
func NewMessage(spec *utils.Specification) (*Message, error) {
	if spec == nil && spec.Elements == nil || spec.Encoding == nil {
		return nil, errors.New("has invalid specification")
	}
	_elements, err := NewDataElements(spec)
	if err != nil {
		return nil, err
	}
	return &Message{
		Mti: &Element{
			Type:   utils.ElementTypeMti,
			Fixed:  true,
			Length: 4,
		},
		Bitmap: &Element{
			Type:   utils.ElementTypeBitmap,
			Fixed:  true,
			Length: 64,
		},
		Elements: _elements,
		Spec:     spec,
	}, nil
}

func (m *Message) Validate() error {
	if m.Mti != nil {
		if err := m.Mti.Validate(); err != nil {
			return err
		}
	}
	if m.Bitmap != nil {
		if err := m.Bitmap.Validate(); err != nil {
			return err
		}
	}
	if m.Elements != nil {
		if err := m.Elements.Validate(); err != nil {
			return err
		}
	}
	return nil
}
