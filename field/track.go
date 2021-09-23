package field

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var _ Field = (*Track)(nil)

type Track struct {
	spec *Spec  `xml:"-" json:"-"`
	data *Track `xml:"-" json:"-"`

	// Available for Track 1
	FixedLength bool `xml:"-" json:"-"`

	// Available track number: 1,2,3
	Number TrackNumber `xml:"number,omitempty" json:"number,omitempty"`

	// Available for Track 1, Track 3
	FormatCode string `xml:"FormatCode,omitempty" json:"format_code,omitempty"`

	// Available for Track 1, Track 2, Track 3
	PrimaryAccountNumber string `xml:"PrimaryAccountNumber,omitempty" json:"primary_account_number,omitempty"`

	// Available for Track 1
	Name string `xml:"Name,omitempty" json:"name,omitempty"`

	// Available for Track 1, Track 2
	ExpirationDate *time.Time `xml:"ExpirationDate,omitempty" json:"expiration_date,omitempty"`

	// Available for Track 1, Track 2
	ServiceCode string `xml:"ServiceCode,omitempty" json:"service_code,omitempty"`

	// Available for Track 1, Track 2, Track 3
	//  If track 3, the field describe security data + additional data
	DiscretionaryData string `xml:"DiscretionaryData,omitempty" json:"discretionary_data,omitempty"`
}

type TrackNumber int

func (v TrackNumber) Valid() bool {
	if v == Track1 || v == Track2 || v == Track3 {
		return true
	}
	return false
}

const (
	Track1 TrackNumber = 1
	Track2 TrackNumber = 2
	Track3 TrackNumber = 3

	expiryDateFormat = "0601"

	track1Format = `%s%s^%s^%s%s%s`
	track2Format = `%s=%s%s%s`
	track3Format = `%s%s=%s`
)

var (
	track1Regex = regexp.MustCompile(`^([A-Z]{1})([0-9]{1,19})\^([^\^]{2,26})\^([0-9]{4}|\^)([0-9]{3}|\^)([^\?]+)$`)
	track2Regex = regexp.MustCompile(`^([0-9]{1,19})\=([0-9]{4}|\=)([0-9]{3}|\=)([^\?]+)$`)
	track3Regex = regexp.MustCompile(`^([0-9]{2})([0-9]{1,19})\=([^\?]+)$`)
)

func NewTrack(spec *Spec, number TrackNumber) (*Track, error) {
	if !number.Valid() {
		return nil, errors.New("invalid track number")
	}
	return &Track{
		spec:   spec,
		Number: number,
	}, nil
}

func NewTrackValue(val []byte, number TrackNumber, fixedLength bool) (*Track, error) {
	if !number.Valid() {
		return nil, errors.New("invalid track number")
	}

	track := &Track{
		Number:      number,
		FixedLength: fixedLength,
	}
	err := track.parse(val)
	if err != nil {
		return nil, errors.New("invalid track data")
	}
	return track, nil
}

func (f *Track) Spec() *Spec {
	return f.spec
}

func (f *Track) SetSpec(spec *Spec) {
	f.spec = spec
}

func (f *Track) SetBytes(b []byte) error {
	return f.parse(b)
}

func (f *Track) Bytes() ([]byte, error) {
	return f.serialize()
}

func (f *Track) String() (string, error) {
	b, err := f.serialize()
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %v", err)
	}
	return string(b), nil
}

func (f *Track) Pack() ([]byte, error) {
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
func (f *Track) Unpack(data []byte) (int, error) {
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

	if f.data != nil {
		*(f.data) = *f
	}

	return read + prefBytes, nil
}

func (f *Track) SetData(data interface{}) error {
	if data == nil {
		return nil
	}

	track, ok := data.(*Track)
	if !ok {
		return fmt.Errorf("data does not match required *Track type")
	}

	if !track.Number.Valid() {
		return fmt.Errorf("contains invalid track number")
	}

	f.data = track
	f.FixedLength = track.FixedLength
	f.Number = track.Number
	f.FormatCode = track.FormatCode
	f.PrimaryAccountNumber = track.PrimaryAccountNumber
	f.Name = track.Name
	f.ExpirationDate = track.ExpirationDate
	f.ServiceCode = track.ServiceCode
	f.DiscretionaryData = track.DiscretionaryData

	return nil
}

func (f *Track) parse(raw []byte) error {
	if raw == nil {
		return errors.New("invalid track data")
	}

	switch f.Number {
	case Track1:
		return f.parseForTrack1(raw)
	case Track2:
		return f.parseForTrack2(raw)
	case Track3:
		return f.parseForTrack3(raw)
	}

	return errors.New("invalid track number")
}

func (f *Track) parseForTrack1(b []byte) error {
	matches := track1Regex.FindStringSubmatch(string(b))
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

	return nil
}

func (f *Track) parseForTrack2(b []byte) error {
	matches := track2Regex.FindStringSubmatch(string(b))
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

	return nil
}

func (f *Track) parseForTrack3(b []byte) error {
	matches := track3Regex.FindStringSubmatch(string(b))
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

	return nil
}

func (f *Track) serialize() ([]byte, error) {

	var raw string

	switch f.Number {
	case Track1:
		name := f.Name
		if len(f.Name) > 1 && f.FixedLength {
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
		raw = fmt.Sprintf(track1Format, f.FormatCode, f.PrimaryAccountNumber, name, expired, code, f.DiscretionaryData)
	case Track2:
		expired := "^"
		if f.ExpirationDate != nil {
			expired = f.ExpirationDate.Format(expiryDateFormat)
		}
		code := "^"
		if len(f.ServiceCode) > 0 {
			code = f.ServiceCode
		}
		raw = fmt.Sprintf(track2Format, f.PrimaryAccountNumber, expired, code, f.DiscretionaryData)
	case Track3:
		raw = fmt.Sprintf(track3Format, f.FormatCode, f.PrimaryAccountNumber, f.DiscretionaryData)
	default:
		return nil, errors.New("unsupported track number")
	}

	return []byte(raw), nil
}
