package field

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*HexNumeric)(nil)
var _ json.Marshaler = (*HexNumeric)(nil)
var _ json.Unmarshaler = (*HexNumeric)(nil)

// HexNumeric is a Hex field that should be represented numerically.
// For example, given the raw bytes []byte{0x12, 0x34}, decoding them as
// a HexNumeric field will result in the value int64(1234).
// It's useful for fields that will likely be used for further arithmetic operations.
// It should only be used with binary encoding.
type HexNumeric struct {
	value int64
	spec  *Spec
}

// NewHexNumeric creates a new instance of the *HexNumeric struct,
// validates and sets its Spec before returning it.
// Refer to SetSpec() for more information on Spec validation.
func NewHexNumeric(spec *Spec) *HexNumeric {
	return &HexNumeric{
		spec: spec,
	}
}

// NewHexNumericValue creates a new HexNumeric field with the given value. The value is
// converted from hex numeric to bytes before packing, so we don't validate that val is
// a valid hex numeric here.
func NewHexNumericValue(val int64) *HexNumeric {
	return &HexNumeric{
		value: val,
	}
}

// Spec returns the receiver's spec.
func (f *HexNumeric) Spec() *Spec {
	return f.spec
}

// SetSpec sets the spec of *HexNumeric.
func (f *HexNumeric) SetSpec(spec *Spec) {
	f.spec = spec
}

// SetBytes sets hex bytes as int64 value.
// For example for []byte{0x01, 0x23} this will set int64 123.
func (f *HexNumeric) SetBytes(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	i, err := strconv.ParseInt(fmt.Sprintf("%X", b), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert into number: %w", err)
	}
	f.value = i
	return nil
}

// Bytes returns the hex bytes representation of the int64 value.
//
//	For example for int64 value 123 this will return []byte{0x01, 0x23}.
func (f *HexNumeric) Bytes() ([]byte, error) {
	if f == nil {
		return nil, nil
	}

	str := fmt.Sprintf("%d", f.value)
	if len(str)%2 != 0 {
		// add leading zero to avoid odd length hex string error
		str = fmt.Sprintf("0%d", f.value)
	}

	bytes, err := hex.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("failed to convert into bytes: %w", err)
	}

	return bytes, nil
}

// String returns the string representation of the HexNumeric value.
// For example if the value is int64 100, this will return "100".
func (f *HexNumeric) String() (string, error) {
	if f == nil {
		return "", nil
	}
	return strconv.FormatInt(f.value, 10), nil
}

// Value returns the int64 value of the HexNumeric.
func (f *HexNumeric) Value() int64 {
	if f == nil {
		return 0
	}
	return f.value
}

// SetValue sets the int64 value on the field.
func (f *HexNumeric) SetValue(v int64) {
	f.value = v
}

// Pack serialises data held by the receiver
// into bytes and returns an error on failure.
func (f *HexNumeric) Pack() ([]byte, error) {
	data, err := f.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bytes: %w", err)
	}

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

// Unpack takes in a byte array and deserializes it into the receiver's
// groups of subfields. An offset (unit depends on encoding and prefix values) is
// returned on success. A non-nil error is returned on failure.
func (f *HexNumeric) Unpack(data []byte) (int, error) {
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

	f.SetBytes(raw)

	return read + prefBytes, nil
}

// SetData is deprecated. Use Marshal instead.
func (f *HexNumeric) SetData(data interface{}) error {
	return f.Marshal(data)
}

// Unmarshal unmarshals FROM HexNumeric TO another data structure.
// v must be a pointer to a reflect.String, reflect.Int64, string, int64, or
// another HexNumeric struct.
func (f *HexNumeric) Unmarshal(v interface{}) error {
	switch val := v.(type) {
	case reflect.Value:
		if !val.CanSet() {
			return fmt.Errorf("cannot set reflect.Value of type %s", val.Kind())
		}

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			s, err := f.String()
			if err != nil {
				return fmt.Errorf("cannot retrieve string value: %w", err)
			}
			val.SetString(s)
		case reflect.Int64:
			val.SetInt(f.value)
		default:
			return fmt.Errorf("unsupported reflect.Value type: %s", val.Kind())
		}
	case *string:
		s, err := f.String()
		if err != nil {
			return fmt.Errorf("cannot retrieve string value: %w", err)
		}
		*val = s
	case *int64:
		*val = f.value
	case *HexNumeric:
		val.value = f.value
	default:
		return fmt.Errorf("unsupported type: expected *HexNumeric, *string, *int, or reflect.Value, got %T", v)
	}

	return nil
}

// Marshal marshals FROM some data structure (v) TO HexNumeric.
// v can be reflect.String, reflect.Int64, string, int64, or
// another HexNumeric struct.
func (f *HexNumeric) Marshal(v interface{}) error {
	if v == nil || reflect.ValueOf(v).IsZero() {
		f.value = 0
		return nil
	}

	switch v := v.(type) {
	case *HexNumeric:
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
		return fmt.Errorf("data does not match required *HexNumeric or (int, *int, string, *string) type")
	}

	return nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
// It marshals from the receiver's subfields to a JSON object.
func (f *HexNumeric) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(f.value)
	if err != nil {
		return nil, utils.NewSafeError(err, "failed to JSON marshal int to bytes")
	}
	return bytes, nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// An error is thrown if the JSON consists of a subfield that has not
// been defined in the spec.
// UnmarshalJSON (for the sake of clarity) unmarshals FROM bytes TO HexNumeric.
func (f *HexNumeric) UnmarshalJSON(b []byte) error {
	var v int
	err := json.Unmarshal(b, &v)
	if err != nil {
		return utils.NewSafeError(err, "failed to JSON unmarshal bytes to int")
	}
	f.SetValue(int64(v))
	return nil
}
