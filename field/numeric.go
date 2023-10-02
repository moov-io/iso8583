package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*Numeric)(nil)
var _ json.Marshaler = (*Numeric)(nil)
var _ json.Unmarshaler = (*Numeric)(nil)

type Numeric struct {
	value int
	spec  *Spec
}

func NewNumeric(spec *Spec) *Numeric {
	return &Numeric{
		spec: spec,
	}
}

func NewNumericValue(val int) *Numeric {
	return &Numeric{
		value: val,
	}
}

func (f *Numeric) Spec() *Spec {
	return f.spec
}

func (f *Numeric) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Numeric) SetBytes(b []byte) error {
	if len(b) == 0 {
		// for a length 0 raw, string(raw) would become "" which makes Atoi return an error
		// however for example "0000" (value 0 left-padded with '0') should have 0 as output, not an error
		// so if the length of raw is 0, set f.value to 0 instead of parsing the raw
		f.value = 0
	} else {
		// otherwise parse the raw to an int
		val, err := strconv.Atoi(string(b))
		if err != nil {
			return utils.NewSafeError(err, "failed to convert into number")
		}
		f.value = val
	}

	return nil
}

func (f *Numeric) Bytes() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return []byte(strconv.Itoa(f.value)), nil
}

func (f *Numeric) String() (string, error) {
	if f == nil {
		return "", nil
	}
	return strconv.Itoa(f.value), nil
}

func (f *Numeric) Value() int {
	if f == nil {
		return 0
	}
	return f.value
}

func (f *Numeric) SetValue(v int) {
	f.value = v
}

func (f *Numeric) Pack() ([]byte, error) {
	data := []byte(strconv.Itoa(f.value))

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

// returns number of bytes was read
func (f *Numeric) Unpack(data []byte) (int, error) {
	dataLen, prefBytes, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %w", err)
	}

	raw, read, err := f.spec.Enc.Decode(data[prefBytes:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %w", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	if err := f.SetBytes(raw); err != nil {
		return 0, fmt.Errorf("failed to set bytes: %w", err)
	}

	return read + prefBytes, nil
}

// Deprecated. Use Marshal instead
func (f *Numeric) SetData(data interface{}) error {
	return f.Marshal(data)
}

func (f *Numeric) Unmarshal(v interface{}) error {
	switch val := v.(type) {
	case reflect.Value:
		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			if !val.CanSet() {
				return fmt.Errorf("reflect.Value of the data can not be change")
			}

			str := strconv.Itoa(f.value)
			val.SetString(str)
		case reflect.Int:
			if !val.CanSet() {
				return fmt.Errorf("reflect.Value of the data can not be change")
			}

			val.SetInt(int64(f.value))
		default:
			return fmt.Errorf("data does not match required reflect.Value type")
		}
	case *string:
		str := strconv.Itoa(f.value)
		*val = str
	case *int:
		*val = f.value
	case *Numeric:
		val.value = f.value
	default:
		return fmt.Errorf("data does not match required *Numeric or *int type")
	}

	return nil
}

func (f *Numeric) Marshal(data interface{}) error {
	switch v := data.(type) {
	case *Numeric:
		if v == nil {
			f.value = 0
			return nil
		}
		f.value = v.value
	case int:
		f.value = v
	case *int:
		if v == nil {
			f.value = 0
			return nil
		}
		f.value = *v
	case string:
		if v == "" {
			f.value = 0
			return nil
		}
		val, err := strconv.Atoi(v)
		if err != nil {
			return utils.NewSafeError(err, "failed to convert sting value into number")
		}
		f.value = val
	case *string:
		if v == nil {
			f.value = 0
			return nil
		}

		val, err := strconv.Atoi(*v)
		if err != nil {
			return utils.NewSafeError(err, "failed to convert sting value into number")
		}
		f.value = val
	default:
		return fmt.Errorf("data does not match require *Numeric or (int, *int, string, *string) type")
	}

	return nil
}

func (f *Numeric) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(f.value)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal int to bytes")
	}
	return bytes, nil
}

func (f *Numeric) UnmarshalJSON(b []byte) error {
	var v int
	err := json.Unmarshal(b, &v)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to int")
	}
	return f.SetBytes([]byte(fmt.Sprintf("%d", v)))
}
