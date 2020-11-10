// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"errors"
	"github.com/moov-io/iso8583/pkg/utils"
)

type Iso8583Message interface {
	Bytes() ([]byte, error)
	Load(raw []byte) (int, error)
	Validate() error
}

// message is structure for ISO 8583 message encode and decode
type isoMessage struct {
	Mti          *Element             `xml:"mti,omitempty" json:"mti,omitempty" yaml:"mti,omitempty"`
	Bitmap       *Element             `xml:"bitmap,omitempty" json:"bitmap,omitempty" yaml:"bitmap,omitempty"`
	Elements     *DataElements        `xml:"elements,omitempty" json:"message,elements" yaml:"elements,omitempty"`
	Spec         *utils.Specification `xml:"-" json:"-"`
	SecondBitmap bool                 `xml:"-" json:"-"`
	ThirdBitmap  bool                 `xml:"-" json:"-"`
}

// create data elements of message with specification
func NewMessage(spec *utils.Specification) (Iso8583Message, error) {
	if spec == nil && spec.Elements == nil || spec.Encoding == nil {
		return nil, errors.New("has invalid specification")
	}
	_elements, err := NewDataElements(spec)
	if err != nil {
		return nil, err
	}
	return &isoMessage{
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

func (m *isoMessage) Validate() error {
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

func (m *isoMessage) Bytes() ([]byte, error) {
	return nil, nil
}

func (m *isoMessage) Load(raw []byte) (int, error) {
	return 0, nil
}
