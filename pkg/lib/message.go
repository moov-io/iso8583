// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/moov-io/iso8583/pkg/utils"
	"reflect"
)

type Iso8583Message interface {
	Bytes() ([]byte, error)
	Load(raw []byte) (int, error)
	Validate() error
	GetElements() map[int]*Element
	GetMti() *Element
	GetBitmap() *Element
}

// create data elements of message with specification
func NewMessage(spec *utils.Specification) (Iso8583Message, error) {
	elements, err := NewDataElements(spec)
	if err != nil {
		return nil, err
	}
	return &isoMessage{
		mti: &Element{
			Type:     utils.ElementTypeMti,
			Fixed:    true,
			Length:   4,
			Encoding: spec.Encoding.MtiEnc,
		},
		bitmap: &Element{
			Type:     utils.ElementTypeBitmap,
			Fixed:    true,
			Length:   64,
			Encoding: spec.Encoding.BitmapEnc,
		},
		elements: elements,
		spec:     spec,
	}, nil
}

// message instance
// isoMessage is structure for ISO 8583 message encode and decode
type isoMessage struct {
	mti      *Element
	bitmap   *Element
	elements *dataElements
	spec     *utils.Specification
	indexes  []int
}

// isoMessage is structure for marshaling and un-marshaling
type messageJSON struct {
	MTI      *Element      `xml:"mti,omitempty" json:"mti,omitempty" yaml:"mti,omitempty"`
	Bitmap   *Element      `xml:"bitmap,omitempty" json:"bitmap,omitempty" yaml:"bitmap,omitempty"`
	Elements *dataElements `xml:"elements,omitempty" json:"elements,omitempty" yaml:"elements,omitempty"`
}

// Validate check validation of field
func (m *isoMessage) Validate() error {
	if m.mti != nil {
		if err := m.mti.Validate(); err != nil {
			return err
		}
	}
	if m.bitmap != nil {
		if err := m.bitmap.Validate(); err != nil {
			return err
		}
	}
	if m.elements != nil {
		if err := m.elements.Validate(); err != nil {
			return err
		}
	}

	m.generateIndexes()
	if !reflect.DeepEqual(m.elements.Keys(), m.indexes) {
		return errors.New(utils.ErrMisMatchElementsBitmap)
	}
	return nil
}

// Bytes encode field to bytes
func (m *isoMessage) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	if m.mti != nil {
		value, err := m.mti.Bytes()
		if err != nil {
			return nil, err
		}
		buf.Write(value)
	}
	if m.bitmap != nil {
		value, err := m.bitmap.Bytes()
		if err != nil {
			return nil, err
		}
		buf.Write(value)
	}
	if m.elements != nil {
		value, err := m.elements.Bytes()
		if err != nil {
			return nil, err
		}
		buf.Write(value)
	}

	return buf.Bytes(), nil
}

// Load decode field from bytes
func (m *isoMessage) Load(raw []byte) (int, error) {
	if m.mti == nil && m.bitmap == nil {
		return 0, errors.New(utils.ErrNonInitializedMessage)
	}

	start := 0
	read, err := m.mti.Load(raw)
	if err != nil {
		return 0, err
	}
	start += read

	read, err = m.bitmap.Load(raw[start:])
	if err != nil {
		return 0, err
	}
	start += read

	m.generateIndexes()
	for _, index := range m.indexes {
		if index > 2 { // second, third bitmap
			break
		}
		read, err = m.createElement(index, start, raw)
		if err != nil {
			return 0, err
		}
		start += read
		m.generateIndexes()
	}

	for _, index := range m.indexes {
		if index < 3 { // second, third bitmap
			continue
		}
		read, err = m.createElement(index, start, raw)
		if err != nil {
			return 0, err
		}
		start += read
	}

	if start != len(raw) {
		return read, errors.New(utils.ErrBadRaw)
	}

	return start, nil
}

// GetElements return data elements of iso message
func (m *isoMessage) GetElements() map[int]*Element {
	if m.elements == nil {
		return nil
	}
	return m.elements.elements
}

// GetMti return mti of iso message
func (m *isoMessage) GetMti() *Element {
	return m.mti
}

// GetBitmap return bitmap of iso message
func (m *isoMessage) GetBitmap() *Element {
	return m.bitmap
}

// Customize unmarshal of json
func (m *isoMessage) UnmarshalJSON(b []byte) error {
	dummy := messageJSON{
		MTI:      m.mti,
		Bitmap:   m.bitmap,
		Elements: m.elements,
	}
	if err := json.Unmarshal(b, &dummy); err != nil {
		return err
	}

	m.mti = dummy.MTI
	m.bitmap = dummy.Bitmap
	m.elements = dummy.Elements
	m.generateIndexes()
	return nil
}

// Customize marshal of json
func (m *isoMessage) MarshalJSON() ([]byte, error) {
	dummy := messageJSON{
		MTI:      m.mti,
		Bitmap:   m.bitmap,
		Elements: m.elements,
	}
	return json.Marshal(dummy)
}

// Customize unmarshal of xml
func (m *isoMessage) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	dummy := messageJSON{
		MTI:      m.mti,
		Bitmap:   m.bitmap,
		Elements: m.elements,
	}
	if err := decoder.DecodeElement(&dummy, &start); err != nil {
		return err
	}

	m.mti = dummy.MTI
	m.bitmap = dummy.Bitmap
	m.elements = dummy.Elements
	m.generateIndexes()
	return nil
}

// Customize marshal of xml
func (m *isoMessage) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	dummy := messageJSON{
		MTI:      m.mti,
		Bitmap:   m.bitmap,
		Elements: m.elements,
	}
	return encoder.EncodeElement(dummy, start)
}

// private functions ...
func (m *isoMessage) generateIndexes() {
	if m.bitmap == nil {
		return
	}
	m.indexes = utils.BitmapToIndexArray(m.bitmap.String(), 0)
	if utils.IsSecondBitmap(m.bitmap.String()) {
		if m.elements.elements[1] != nil {
			indexes := utils.BitmapToIndexArray(m.elements.elements[1].String(), 64)
			m.indexes = append(m.indexes, indexes...)
		}
		if utils.IsThirdBitmap(m.bitmap.String()) {
			if m.elements.elements[2] != nil &&
				m.elements.elements[2].Type == utils.ElementTypeBinary &&
				m.elements.elements[2].Length == 64 {
				indexes := utils.BitmapToIndexArray(m.elements.elements[2].String(), 128)
				m.indexes = append(m.indexes, indexes...)
			}
		}
	}
}

func (m *isoMessage) createElement(index, start int, raw []byte) (int, error) {
	spec, err := m.spec.Elements.Get(index)
	if err != nil {
		return 0, err
	}
	_type, err := spec.Parse()
	if err != nil {
		return 0, err
	}
	_type.SetEncoding(m.spec.Encoding)
	elm := &Element{}
	elm.setType(_type)

	if start >= len(raw) {
		return 0, errors.New(utils.ErrBadRaw)
	}

	read, err := elm.Load(raw[start:])
	if err != nil {
		return 0, err
	}
	m.elements.elements[index] = elm

	return read, nil
}
