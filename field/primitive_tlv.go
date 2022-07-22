package field

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moov-io/iso8583/encoding"
)

var _ Field = (*PrimitiveTLV)(nil)
var _ json.Marshaler = (*PrimitiveTLV)(nil)
var _ json.Unmarshaler = (*PrimitiveTLV)(nil)

// A data element is the value field (V) of a primitive BER-TLV data object. A data
// element is the smallest data field that receives an identifier (a tag).
//
//  Tag  Length  Value
//  (T)   (L)     (V)
//

type PrimitiveTLV struct {
	Value []byte `json:"value,omitempty"`

	spec *Spec
	data *PrimitiveTLV
}

// NewPrimitiveTLV returns a instance of primitive tlv with spec
func NewPrimitiveTLV(spec *Spec) *PrimitiveTLV {
	return &PrimitiveTLV{
		spec: spec,
	}
}

// NewPrimitiveTLVValue returns a instance of primitive tlv with value (raw bytes)
func NewPrimitiveTLVValue(val []byte) *PrimitiveTLV {
	return &PrimitiveTLV{
		Value: val,
	}
}

// NewPrimitiveTLVHexString returns a instance of primitive tlv with hex string
func NewPrimitiveTLVHexString(val string) *PrimitiveTLV {
	value, err := encoding.BerTLVTag.Encode([]byte(val))

	if err != nil {
		return &PrimitiveTLV{}
	}

	return &PrimitiveTLV{
		Value: value,
	}
}

// Spec returns a specification of tlv field
func (f *PrimitiveTLV) Spec() *Spec {
	return f.spec
}

// SetSpec set specification into a tlv field
func (f *PrimitiveTLV) SetSpec(spec *Spec) {
	f.spec = spec
}

// SetBytes set value of tlv field (only value)
func (f *PrimitiveTLV) SetBytes(b []byte) error {
	f.Value = b
	if f.data != nil {
		*(f.data) = *f
	}
	return nil
}

// Bytes returns raw value of tlv
func (f *PrimitiveTLV) Bytes() ([]byte, error) {
	return f.Value, nil
}

// Bytes returns hex string of value of tlv
func (f *PrimitiveTLV) String() (string, error) {
	return fmt.Sprintf("%X", f.Value), nil
}

// Pack returns encoded bytes for tlv (Tag + Length + Value)
func (f *PrimitiveTLV) Pack() ([]byte, error) {

	if f.spec.Tag == nil || f.spec.Tag.Enc == nil || f.spec.Pref == nil || f.spec.Enc == nil {
		return nil, fmt.Errorf("failed to pack tlv: invalid spec")
	}

	data := f.Value

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packedValue, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	packed := packedValue

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	packed = append(packedLength, packed...)

	tagBytes := []byte(f.spec.Tag.Tag)
	if f.spec.Tag.Pad != nil {
		tagBytes = f.spec.Tag.Pad.Pad(tagBytes, f.spec.Tag.Length)
	}

	tagBytes, err = f.spec.Tag.Enc.Encode(tagBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert subfield Tag \"%v\" to int", tagBytes)
	}

	packed = append(tagBytes, packed...)

	return packed, nil
}

// Unpack do decode tlv field with encoded raw data (Tag + Length + Value)
func (f *PrimitiveTLV) Unpack(data []byte) (int, error) {

	if f.spec.Tag == nil || f.spec.Tag.Enc == nil || f.spec.Pref == nil || f.spec.Enc == nil {
		return 0, fmt.Errorf("failed to unpack tlv: invalid spec")
	}

	offset := 0

	// 1. Read Tag

	tagBytes, read, err := f.spec.Tag.Enc.Decode(data[offset:], f.spec.Tag.Length)
	if err != nil {
		return 0, fmt.Errorf("failed to unpack subfield Tag: %w", err)
	}
	offset += read

	if f.spec.Tag.Pad != nil {
		tagBytes = f.spec.Tag.Pad.Unpad(tagBytes)
	}
	tag := string(tagBytes)

	if tag != f.spec.Tag.Tag {
		return 0, fmt.Errorf("tag mismatch: want to read %s, got %s", f.spec.Tag.Tag, tag)
	}

	// 2. Read Length

	dataLen, read, err := f.spec.Pref.DecodeLength(f.spec.Length, data[offset:])
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}
	offset += read

	// 3. Read Value

	raw, read, err := f.spec.Enc.Decode(data[offset:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %w", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	if err := f.SetBytes(raw); err != nil {
		return 0, fmt.Errorf("failed to set bytes: %w", err)
	}

	return read + offset, nil
}

func (f *PrimitiveTLV) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	tlv, ok := data.(*PrimitiveTLV)
	if !ok {
		return fmt.Errorf("data does not match required *String type")
	}

	f.data = tlv
	if len(tlv.Value) > 0 {
		f.Value = tlv.Value
	}
	return nil
}

func (f *PrimitiveTLV) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *PrimitiveTLV) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	tlv, ok := v.(*PrimitiveTLV)
	if !ok {
		return errors.New("data does not match required *PrimitiveTLV type")
	}

	tlv.Value = f.Value

	return nil
}

func (f *PrimitiveTLV) MarshalJSON() ([]byte, error) {
	str, err := f.String()
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

func (f *PrimitiveTLV) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return fmt.Errorf("failed to JSON unmarshal bytes to string: %w", err)
	}

	hex, err := encoding.ASCIIHexToBytes.Encode([]byte(v))
	if err != nil {
		return fmt.Errorf("failed to convert ASCII Hex string to bytes")
	}
	return f.SetBytes(hex)
}
