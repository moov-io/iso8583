package field

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

var _ Field = (*String)(nil)

type String struct {
	Value string `json:"value"`
	spec  *Spec
	data  *String
}

func NewString(spec *Spec) *String {
	return &String{
		spec: spec,
	}
}

func NewStringValue(val string) *String {
	return &String{
		Value: val,
	}
}

func (f *String) Spec() *Spec {
	return f.spec
}

func (f *String) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *String) SetBytes(b []byte) error {
	f.Value = string(b)
	return nil
}

func (f *String) Bytes() ([]byte, error) {
	return []byte(f.Value), nil
}

func (f *String) String() (string, error) {
	return f.Value, nil
}

func (f *String) Pack() ([]byte, error) {
	var buf bytes.Buffer

	_, err := f.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func (f *String) Unpack(data []byte) (int, error) {
	return f.ReadFrom(bytes.NewReader(data))
}

func (f *String) WriteTo(w io.Writer) (n int, err error) {
	data := []byte(f.Value)

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return 0, fmt.Errorf("failed to encode content: %v", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return 0, fmt.Errorf("failed to encode length: %v", err)
	}

	m, err := w.Write(packedLength)
	if err != nil {
		return m, fmt.Errorf("writing packed length: %v", err)
	}

	n += m

	m, err = w.Write(packed)
	if err != nil {
		return m, fmt.Errorf("writing packed field: %v", err)
	}

	n += m

	return n, nil
}

func (f *String) ReadFrom(r io.Reader) (int, error) {
	dataLen, err := f.spec.Pref.ReadLength(f.spec.Length, r)
	if err != nil {
		return 0, fmt.Errorf("reading length: %v", err)
	}

	// start := f.spec.Pref.Length()
	raw, read, err := f.spec.Enc.DecodeFrom(r, dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %v", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	f.Value = string(raw)

	if f.data != nil {
		*(f.data) = *f
	}

	return read + f.spec.Pref.Length(), nil
}

func (f *String) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	str, ok := data.(*String)
	if !ok {
		return fmt.Errorf("data does not match required *String type")
	}

	f.data = str
	if str.Value != "" {
		f.Value = str.Value
	}
	return nil
}

func (f *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}
