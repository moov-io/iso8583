package field

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var _ Field = (*Track2)(nil)

type Track2 struct {
	PrimaryAccountNumber string     `xml:"PrimaryAccountNumber,omitempty" json:"primary_account_number,omitempty"`
	Separator            string     `xml:"Separator,omitempty" json:"separator,omitempty"`
	ExpirationDate       *time.Time `xml:"ExpirationDate,omitempty" json:"expiration_date,omitempty"`
	ServiceCode          string     `xml:"ServiceCode,omitempty" json:"service_code,omitempty"`
	DiscretionaryData    string     `xml:"DiscretionaryData,omitempty" json:"discretionary_data,omitempty"`

	spec *Spec
	data *Track2
}

const (
	track2Format = `%s%s%s%s%s`

	defaultSeparator = "="
)

var (
	track2Regex = regexp.MustCompile(`^([0-9]{1,19})(=|D)([0-9]{4})([0-9]{3})([^?]+)$`)
)

func NewTrack2(spec *Spec) *Track2 {
	return &Track2{
		spec: spec,
	}
}

func (f *Track2) Spec() *Spec {
	return f.spec
}

func (f *Track2) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Track2) SetBytes(b []byte) error {
	return f.unpack(b)
}

func (f *Track2) Bytes() ([]byte, error) {
	return f.pack()
}

func (f *Track2) String() (string, error) {
	b, err := f.pack()
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %w", err)
	}
	return string(b), nil
}

func (f *Track2) Pack() ([]byte, error) {
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
func (f *Track2) Unpack(data []byte) (int, error) {
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

func (f *Track2) Unmarshal(v interface{}) error {
	if v == nil {
		return nil
	}

	track, ok := v.(*Track2)
	if !ok {
		return fmt.Errorf("data does not match required *Track2 type")
	}

	track.PrimaryAccountNumber = f.PrimaryAccountNumber
	track.Separator = f.Separator
	track.ExpirationDate = f.ExpirationDate
	track.ServiceCode = f.ServiceCode
	track.DiscretionaryData = f.DiscretionaryData

	return nil
}

func (f *Track2) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	track, ok := data.(*Track2)
	if !ok {
		return fmt.Errorf("data does not match required *Track type")
	}

	f.PrimaryAccountNumber = track.PrimaryAccountNumber
	f.Separator = track.Separator
	f.ExpirationDate = track.ExpirationDate
	f.ServiceCode = track.ServiceCode
	f.DiscretionaryData = track.DiscretionaryData

	f.data = track

	return nil
}

func (f *Track2) Marshal(data interface{}) error {
	return f.SetData(data)
}

func (f *Track2) unpack(raw []byte) error {
	if raw == nil || !track2Regex.Match(raw) {
		return errors.New("invalid track data")
	}

	matches := track2Regex.FindStringSubmatch(string(raw))
	for index, val := range matches {
		value := strings.TrimSpace(val)
		if len(value) == 0 {
			continue
		}

		switch index {
		case 1: // Payment card number (PAN)
			f.PrimaryAccountNumber = value
		case 2: // Separator
			f.Separator = value
		case 3: // Expiration Date (ED)
			expiredTime, timeErr := time.Parse(expiryDateFormat, value)
			if timeErr != nil {
				return errors.New("invalid expired time")
			}
			f.ExpirationDate = &expiredTime
		case 4: // Service Code (SC)
			f.ServiceCode = value
		case 5: // Discretionary data (DD)
			f.DiscretionaryData = value
		}
	}

	if f.data != nil {
		*(f.data) = *f
	}

	return nil
}

func (f *Track2) pack() ([]byte, error) {
	expired := "^"
	if f.ExpirationDate != nil {
		expired = f.ExpirationDate.Format(expiryDateFormat)
	}
	code := "^"
	if len(f.ServiceCode) > 0 {
		code = f.ServiceCode
	}
	separator := defaultSeparator
	if f.Separator != "" {
		separator = f.Separator
	}

	raw := fmt.Sprintf(track2Format, f.PrimaryAccountNumber, separator, expired, code, f.DiscretionaryData)
	return []byte(raw), nil
}
