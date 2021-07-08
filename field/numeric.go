package field

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

var _ Field = (*Numeric)(nil)

type Numeric struct {
	Value int `json:"value"`
	spec  *Spec
	data  *Numeric
}

func NewNumeric(spec *Spec) *Numeric {
	return &Numeric{
		spec: spec,
	}
}

func NewNumericValue(val int) *Numeric {
	return &Numeric{
		Value: val,
	}
}

func (f *Numeric) Spec() *Spec {
	return f.spec
}

func (f *Numeric) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Numeric) SetBytes(b []byte) error {
	val, err := strconv.Atoi(string(b))
	if err == nil {
		f.Value = val
	}
	return err
}

func (f *Numeric) Bytes() ([]byte, error) {
	return []byte(strconv.Itoa(f.Value)), nil
}

func (f *Numeric) String() (string, error) {
	return strconv.Itoa(f.Value), nil
}

func (f *Numeric) Pack() ([]byte, error) {
	data := []byte(strconv.Itoa(f.Value))

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %v", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %v", err)
	}

	return append(packedLength, packed...), nil
}

// returns number of bytes was read
func (f *Numeric) ReadFrom(r io.Reader) (int, error) {
	dataLen, err := f.spec.Pref.ReadLength(f.spec.Length, r)
	if err != nil {
		return 0, fmt.Errorf("reading length: %v", err)
	}

	raw, read, err := f.spec.Enc.DecodeFrom(r, dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %v", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	if len(raw) == 0 {
		// for a length 0 raw, string(raw) would become "" which makes Atoi return an error
		// however for example "0000" (value 0 left-padded with '0') should have 0 as output, not an error
		// so if the length of raw is 0, set f.Value to 0 instead of parsing the raw
		f.Value = 0
	} else {
		// otherwise parse the raw to an int
		f.Value, err = strconv.Atoi(string(raw))
		if err != nil {
			return 0, fmt.Errorf("failed to convert into number: %v", err)
		}
	}

	if f.data != nil {
		*(f.data) = *f
	}

	return read + f.spec.Pref.Length(), nil
}

func (f *Numeric) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	num, ok := data.(*Numeric)
	if !ok {
		return fmt.Errorf("data does not match required *Numeric type")
	}

	f.data = num
	if num.Value != 0 {
		f.Value = num.Value
	}
	return nil
}

func (f *Numeric) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}
