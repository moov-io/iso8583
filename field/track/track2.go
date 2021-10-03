package track

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/moov-io/iso8583/field"
)

var _ field.Field = (*Track2)(nil)

type Track2 struct {
	spec                 *field.Spec `xml:"-" json:"-"`
	PrimaryAccountNumber string      `xml:"PrimaryAccountNumber,omitempty" json:"primary_account_number,omitempty"`
	ExpirationDate       *time.Time  `xml:"ExpirationDate,omitempty" json:"expiration_date,omitempty"`
	ServiceCode          string      `xml:"ServiceCode,omitempty" json:"service_code,omitempty"`
	DiscretionaryData    string      `xml:"DiscretionaryData,omitempty" json:"discretionary_data,omitempty"`

	data *Track2
}

const (
	track2Format = `%s=%s%s%s`
)

var (
	track2Regex = regexp.MustCompile(`^([0-9]{1,19})\=([0-9]{4}|\=)([0-9]{3}|\=)([^\?]+)$`)
)

func NewTrack2(spec *field.Spec) *Track2 {
	return &Track2{
		spec: spec,
	}
}

func NewTrack2Value(val []byte) (*Track2, error) {
	track := &Track2{}
	err := track.parse(val)
	if err != nil {
		return nil, errors.New("invalid track data")
	}
	return track, nil
}

func (f *Track2) Spec() *field.Spec {
	return f.spec
}

func (f *Track2) SetSpec(spec *field.Spec) {
	f.spec = spec
}

func (f *Track2) SetBytes(b []byte) error {
	return f.parse(b)
}

func (f *Track2) Bytes() ([]byte, error) {
	return f.serialize()
}

func (f *Track2) String() (string, error) {
	b, err := f.serialize()
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %v", err)
	}
	return string(b), nil
}

func (f *Track2) Pack() ([]byte, error) {
	data, err := f.serialize()
	if err != nil {
		return nil, err
	}

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
func (f *Track2) Unpack(data []byte) (int, error) {
	dataLen, prefBytes, err := f.spec.Pref.DecodeLength(f.spec.Length, data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode length: %v", err)
	}

	raw, read, err := f.spec.Enc.Decode(data[prefBytes:], dataLen)
	if err != nil {
		return 0, fmt.Errorf("failed to decode content: %v", err)
	}

	if f.spec.Pad != nil {
		raw = f.spec.Pad.Unpad(raw)
	}

	if len(raw) > 0 {
		err = f.parse(raw)
		if err != nil {
			return 0, err
		}
	}
	return read + prefBytes, nil
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
	f.ExpirationDate = track.ExpirationDate
	f.ServiceCode = track.ServiceCode
	f.DiscretionaryData = track.DiscretionaryData

	f.data = track

	return nil
}

func (f *Track2) parse(raw []byte) error {
	if raw == nil || !track2Regex.Match(raw) {
		return errors.New("invalid track data")
	}

	matches := track2Regex.FindStringSubmatch(string(raw))
	for index, val := range matches {
		value := strings.TrimSpace(val)
		if len(value) == 0 || value == "=" {
			continue
		}

		switch index {
		case 1: // Payment card number (PAN)
			f.PrimaryAccountNumber = value
		case 2: // Expiration Date (ED)
			expiredTime, timeErr := time.Parse(expiryDateFormat, value)
			if timeErr != nil {
				return errors.New("invalid expired time")
			}
			f.ExpirationDate = &expiredTime
		case 3: // Service Code (SC)
			f.ServiceCode = value
		case 4: // Discretionary data (DD)
			f.DiscretionaryData = value
		}
	}

	if f.data != nil {
		*(f.data) = *f
	}

	return nil
}

func (f *Track2) serialize() ([]byte, error) {

	expired := "^"
	if f.ExpirationDate != nil {
		expired = f.ExpirationDate.Format(expiryDateFormat)
	}
	code := "^"
	if len(f.ServiceCode) > 0 {
		code = f.ServiceCode
	}

	raw := fmt.Sprintf(track2Format, f.PrimaryAccountNumber, expired, code, f.DiscretionaryData)
	return []byte(raw), nil
}
