package field

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var _ Field = (*Track3)(nil)

type Track3 struct {
	FormatCode           string `json:"format_code,omitempty"`
	PrimaryAccountNumber string `json:"primary_account_number,omitempty"`
	DiscretionaryData    string `json:"discretionary_data,omitempty"`

	spec *Spec
	data *Track3
}

const (
	track3Format = `%s%s=%s`
)

var (
	track3Regex = regexp.MustCompile(`^([0-9]{2})([0-9]{1,19})\=([^\?]+)$`)
)

func NewTrack3(spec *Spec) *Track3 {
	return &Track3{
		spec: spec,
	}
}

func (f *Track3) Spec() *Spec {
	return f.spec
}

func (f *Track3) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Track3) SetBytes(b []byte) error {
	if err := f.unpack(b); err != nil {
		return nil
	}
	return nil
}

func (f *Track3) Bytes() ([]byte, error) {
	return f.pack()
}

func (f *Track3) String() (string, error) {
	b, err := f.pack()
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %w", err)
	}
	return string(b), nil
}

func (f *Track3) Pack() ([]byte, error) {
	data, err := f.pack()
	if err != nil {
		return nil, err
	}

	if f.spec.Pad != nil {
		data = f.spec.Pad.Pad(data, f.spec.Length)
	}

	packed, err := f.spec.Enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode content: %w", err)
	}

	packedLength, err := f.spec.Pref.EncodeLength(f.spec.Length, len(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to encode length: %w", err)
	}

	return append(packedLength, packed...), nil
}

// returns number of bytes was read
func (f *Track3) Unpack(data []byte) (int, error) {
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

	if len(raw) > 0 {
		err = f.unpack(raw)
		if err != nil {
			return 0, err
		}
	}

	return read + prefBytes, nil
}

func (f *Track3) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	track, ok := v.(*Track3)
	if !ok {
		return fmt.Errorf("data does not match required *Track3 type")
	}

	track.PrimaryAccountNumber = f.PrimaryAccountNumber
	track.FormatCode = f.FormatCode
	track.DiscretionaryData = f.DiscretionaryData

	return nil
}

func (f *Track3) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	track, ok := data.(*Track3)
	if !ok {
		return fmt.Errorf("data does not match required *Track type")
	}

	f.FormatCode = track.FormatCode
	f.PrimaryAccountNumber = track.PrimaryAccountNumber
	f.DiscretionaryData = track.DiscretionaryData

	f.data = track

	return nil
}

func (f *Track3) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *Track3) unpack(raw []byte) error {
	if raw == nil || !track3Regex.Match(raw) {
		return errors.New("invalid track data")
	}

	matches := track3Regex.FindStringSubmatch(string(raw))
	for index, val := range matches {
		value := strings.TrimSpace(val)
		if len(value) == 0 || value == "=" {
			continue
		}

		switch index {
		case 1: // Format Code
			f.FormatCode = value
		case 2: // Payment card number (PAN)
			f.PrimaryAccountNumber = value
		case 3: // Security Data + Additional Data
			f.DiscretionaryData = value
		}
	}

	if f.data != nil {
		*(f.data) = *f
	}

	return nil
}

func (f *Track3) pack() ([]byte, error) {
	raw := fmt.Sprintf(track3Format, f.FormatCode, f.PrimaryAccountNumber, f.DiscretionaryData)
	return []byte(raw), nil
}
