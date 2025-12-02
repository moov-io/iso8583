package field

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/moov-io/iso8583/utils"
)

var (
	_ Field            = (*Numeric)(nil)
	_ json.Marshaler   = (*Numeric)(nil)
	_ json.Unmarshaler = (*Numeric)(nil)
)

type Numeric struct {
	value int64
	spec  *Spec
}

func NewNumeric(spec *Spec) *Numeric {
	return &Numeric{
		spec: spec,
	}
}

func NewNumericValue(val int64) *Numeric {
	return &Numeric{
		value: val,
	}
}

// NewInstance creates a new instance of the Numeric field with the same Spec.
func (f *Numeric) NewInstance() Field {
	return &Numeric{
		spec: f.spec,
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
		val, err := strconv.ParseInt(string(b), 10, 64)
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
	return []byte(strconv.FormatInt(f.value, 10)), nil
}

func (f *Numeric) String() (string, error) {
	if f == nil {
		return "", nil
	}
	return strconv.FormatInt(f.value, 10), nil
}

func (f *Numeric) Value() int64 {
	if f == nil {
		return 0
	}
	return f.value
}

func (f *Numeric) SetValue(v int64) {
	f.value = v
}

func (f *Numeric) Pack() ([]byte, error) {
	data := []byte(strconv.FormatInt(f.value, 10))

	packer := f.spec.getPacker()

	return packer.Pack(data, f.spec)
}

// returns number of bytes was read
func (f *Numeric) Unpack(data []byte) (int, error) {
	unpacker := f.spec.getUnpacker()

	raw, bytesRead, err := unpacker.Unpack(data, f.spec)
	if err != nil {
		return 0, err
	}

	if err := f.SetBytes(raw); err != nil {
		return 0, fmt.Errorf("failed to set bytes: %w", err)
	}

	return bytesRead, nil
}

// Deprecated. Use Marshal instead
func (f *Numeric) SetData(data interface{}) error {
	return f.Marshal(data)
}

func (f *Numeric) Unmarshal(v interface{}) error {
	switch val := v.(type) {
	case reflect.Value:
		if !val.CanSet() {
			return fmt.Errorf("cannot set reflect.Value of type %s", val.Kind())
		}

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			str := strconv.FormatInt(f.value, 10)
			val.SetString(str)
		case reflect.Int64:
			val.SetInt(int64(f.value))
		default:
			return fmt.Errorf("unsupported reflect.Value type: %s", val.Kind())
		}
	case *string:
		str := strconv.FormatInt(f.value, 10)
		*val = str
	case *int64:
		*val = f.value
	case *Numeric:
		val.value = f.value
	default:
		return fmt.Errorf("unsupported type: expected *Numeric, *int, or reflect.Value, got %T", v)
	}

	return nil
}

func (f *Numeric) Marshal(v any) error {
	if v == nil || reflect.ValueOf(v).IsZero() {
		f.value = 0
		return nil
	}

	switch v := v.(type) {
	case *Numeric:
		f.value = v.value
	case int64:
		f.value = v
	case *int64:
		f.value = *v
	case string:
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return utils.NewSafeError(err, "failed to convert sting value into number")
		}
		f.value = val
	case *string:
		val, err := strconv.ParseInt(*v, 10, 64)
		if err != nil {
			return utils.NewSafeError(err, "failed to convert sting value into number")
		}
		f.value = val
	default:
		return fmt.Errorf("data does not match require *Numeric or (int64, *int64, string, *string) type")
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
