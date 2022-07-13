package field

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var _ Field = (*Track1)(nil)

type Track1 struct {
	FixedLength          bool       `json:"fixed_length,omitempty"`
	FormatCode           string     `json:"format_code,omitempty"`
	PrimaryAccountNumber string     `json:"primary_account_number,omitempty"`
	Name                 string     `json:"name,omitempty"`
	ExpirationDate       *time.Time `json:"expiration_date,omitempty"`
	ServiceCode          string     `json:"service_code,omitempty"`
	DiscretionaryData    string     `json:"discretionary_data,omitempty"`

	spec *Spec
	data *Track1
}

const (
	expiryDateFormat = "0601"
	track1Format     = `%s%s^%s^%s%s%s`
)

var (
	track1Regex = regexp.MustCompile(`^([A-Z]{1})([0-9]{1,19})\^([^\^]{2,26})\^([0-9]{4}|\^)([0-9]{3}|\^)([^\?]+)$`)
)

func NewTrack1(spec *Spec) *Track1 {
	return &Track1{
		spec: spec,
	}
}

func (f *Track1) Spec() *Spec {
	return f.spec
}

func (f *Track1) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Track1) SetBytes(b []byte) error {
	return f.unpack(b)
}

func (f *Track1) Bytes() ([]byte, error) {
	return f.pack()
}

func (f *Track1) String() (string, error) {
	b, err := f.pack()
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %w", err)
	}
	return string(b), nil
}

func (f *Track1) Pack() ([]byte, error) {
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
func (f *Track1) Unpack(data []byte) (int, error) {
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

func (f *Track1) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	track, ok := v.(*Track1)
	if !ok {
		return fmt.Errorf("data does not match required *Track1 type")
	}

	track.FixedLength = f.FixedLength
	track.FormatCode = f.FormatCode
	track.PrimaryAccountNumber = f.PrimaryAccountNumber
	track.Name = f.Name
	track.ExpirationDate = f.ExpirationDate
	track.ServiceCode = f.ServiceCode
	track.DiscretionaryData = f.DiscretionaryData

	return nil
}

func (f *Track1) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	track, ok := data.(*Track1)
	if !ok {
		return fmt.Errorf("data does not match required *Track type")
	}

	f.FixedLength = track.FixedLength
	f.FormatCode = track.FormatCode
	f.PrimaryAccountNumber = track.PrimaryAccountNumber
	f.Name = track.Name
	f.ExpirationDate = track.ExpirationDate
	f.ServiceCode = track.ServiceCode
	f.DiscretionaryData = track.DiscretionaryData

	f.data = track

	return nil
}

func (f *Track1) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *Track1) unpack(raw []byte) error {
	if raw == nil || !track1Regex.Match(raw) {
		return errors.New("invalid track data")
	}

	matches := track1Regex.FindStringSubmatch(string(raw))
	for index, val := range matches {
		value := strings.TrimSpace(val)
		if len(value) == 0 || value == "^" {
			continue
		}

		switch index {
		case 1: // Format Code
			f.FormatCode = value
		case 2: // Payment card number (PAN)
			f.PrimaryAccountNumber = value
		case 3: // Name (NM)
			f.Name = value
		case 4: // Expiration Date (ED)
			expiredTime, timeErr := time.Parse(expiryDateFormat, value)
			if timeErr != nil {
				return errors.New("invalid expired time")
			}
			f.ExpirationDate = &expiredTime
		case 5: // Service Code (SC)
			f.ServiceCode = value
		case 6: // Discretionary data (DD)
			f.DiscretionaryData = value
		}
	}

	if f.data != nil {
		*(f.data) = *f
	}

	return nil
}

func (f *Track1) pack() ([]byte, error) {
	name := f.Name
	if len(f.Name) > 1 && f.FixedLength {
		// limit Name to 26 runes and padd with spaces on the right
		name = fmt.Sprintf("%-26.26s", f.Name)
	}
	expired := "^"
	if f.ExpirationDate != nil {
		expired = f.ExpirationDate.Format(expiryDateFormat)
	}
	code := "^"
	if len(f.ServiceCode) > 0 {
		code = f.ServiceCode
	}

	raw := fmt.Sprintf(track1Format, f.FormatCode, f.PrimaryAccountNumber, name, expired, code, f.DiscretionaryData)
	return []byte(raw), nil
}
