package field

import (
	"errors"
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"
	"time"

	"github.com/moov-io/iso8583/utils"
)

var _ Field = (*Track)(nil)

type Track struct {
	spec *Spec  `xml:"-" json:"-"`
	data *Track `xml:"-" json:"-"`

	// Available for Track 1
	FixedLength bool `xml:"-" json:"-"`

	// Available versions 1,2,3
	Version TrackVersion `xml:"Version,omitempty" json:"version,omitempty"`

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

type TrackVersion int

func (v TrackVersion) Valid() bool {
	if v == VersionFirst || v == VersionSecond || v == VersionThird {
		return true
	}
	return false
}

const (
	VersionFirst  = 1
	VersionSecond = 2
	VersionThird  = 3
)

var (
	expiryDateFormat   = "0601"
	trackFirstRegex    = regexp.MustCompile(`^([A-Z]{1})([0-9]{1,19})\^([^\^]{2,26})\^([0-9]{4}|\^)([0-9]{3}|\^)([^\?]+)$`)
	trackSecondRegex   = regexp.MustCompile(`^([0-9]{1,19})\=([0-9]{4}|\=)([0-9]{3}|\=)([^\?]+)$`)
	trackThirdRegex    = regexp.MustCompile(`^([0-9]{2})([0-9]{1,19})\=([^\?]+)$`)
	trackFirstPattern  = `^([A-Z]{1})([0-9]{1,19})\^([^\^]{2,26})\^([0-9]{4}|\^)([0-9]{3}|\^)([^\?]+)$`
	trackSecondPattern = `^([0-9]{1,19})\=([0-9]{4}|\=)([0-9]{3}|\=)([^\?]+)$`
	trackThirdPattern  = `^([0-9]{2})([0-9]{1,19})\=([^\?]+)$`
)

func NewTrack(spec *Spec, version TrackVersion) (*Track, error) {
	if !version.Valid() {
		return nil, errors.New("invalid track version")
	}
	return &Track{
		spec:    spec,
		Version: version,
	}, nil
}

func NewTrack1Value(val []byte, version TrackVersion, isFixed bool) (*Track, error) {
	if !version.Valid() {
		return nil, errors.New("invalid track version")
	}

	track := &Track{
		Version:     VersionFirst,
		FixedLength: isFixed,
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
		return fmt.Errorf("data does not match required *Numeric type")
	}

	f.data = track
	if track.Version == VersionFirst || track.Version == VersionSecond || track.Version == VersionThird {
		f.FixedLength = track.FixedLength
		f.Version = track.Version
		f.FormatCode = track.FormatCode
		f.PrimaryAccountNumber = track.PrimaryAccountNumber
		f.Name = track.Name
		f.ExpirationDate = track.ExpirationDate
		f.ServiceCode = track.ServiceCode
		f.DiscretionaryData = track.DiscretionaryData
	}
	return nil
}

var cardTypes = map[string]string{
	"Visa":             `^4[0-9]{12}(?:[0-9]{3})?$`,
	"MasterCard":       `^5[1-5][0-9]{14}$`,
	"American Express": `^3[47][0-9]{13}$`,
	"Diners Club":      `^3(?:0[0-5]|[68][0-9])[0-9]{11}$`,
	"Discover":         `^6(?:011|5[0-9]{2})[0-9]{12}$`,
	"JCB":              `^(?:2131|1800|35\d{3})\d{11}$`,
}

func (f *Track) GetCardType() string {
	for key, pattern := range cardTypes {
		r, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		if !r.MatchString(f.PrimaryAccountNumber) {
			return key
		}
	}
	return "Unknown"
}

func (f *Track) parse(raw []byte) error {
	if raw == nil {
		return errors.New("invalid track data")
	}

	switch f.Version {
	case VersionFirst:
		return f.parseForVersionFirst(raw)
	case VersionSecond:
		return f.parseForVersionSecond(raw)
	case VersionThird:
		return f.parseForVersionThird(raw)
	}

	return errors.New("invalid track version")
}

func (f *Track) parseForVersionFirst(b []byte) error {
	matches := trackFirstRegex.FindStringSubmatch(string(b))
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
			if timeErr == nil {
				f.ExpirationDate = &expiredTime
			}
		case 5: // Service Code (SC)
			f.ServiceCode = value
		case 6: // Discretionary data (DD)
			f.DiscretionaryData = value
		}
	}

	return nil
}

func (f *Track) parseForVersionSecond(b []byte) error {
	matches := trackSecondRegex.FindStringSubmatch(string(b))
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
			if timeErr == nil {
				f.ExpirationDate = &expiredTime
			}
		case 3: // Service Code (SC)
			f.ServiceCode = value
		case 4: // Discretionary data (DD)
			f.DiscretionaryData = value
		}
	}

	return nil
}

func (f *Track) parseForVersionThird(b []byte) error {
	matches := trackThirdRegex.FindStringSubmatch(string(b))
	for index, val := range matches {
		value := strings.TrimSpace(val)
		if len(value) == 0 || value == "=" {
			continue
		}

		switch index {
		case 1: // Payment card number (PAN)
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

	var pattern string
	var handler func(index int, name string, group *syntax.Regexp, generator utils.Generator, args *utils.GeneratorArgs) string

	switch f.Version {
	case VersionFirst:
		pattern = trackFirstPattern
		handler = func(index int, name string, group *syntax.Regexp, generator utils.Generator, args *utils.GeneratorArgs) string {
			var raw string
			switch index {
			case 0:
				raw = f.FormatCode
			case 1:
				raw = f.PrimaryAccountNumber
			case 2:
				if len(f.Name) > 0 && f.FixedLength {
					raw = fmt.Sprintf("%-26.26s", f.Name)
				} else {
					raw = f.Name
				}
			case 3:
				if f.ExpirationDate != nil {
					raw = f.ExpirationDate.Format(expiryDateFormat)
				}
			case 4:
				raw = f.ServiceCode
			case 5:
				raw = f.DiscretionaryData
			}

			if len(raw) == 0 {
				return `^`
			}
			return raw
		}
	case VersionSecond:
		pattern = trackSecondPattern
		handler = func(index int, name string, group *syntax.Regexp, generator utils.Generator, args *utils.GeneratorArgs) string {
			var raw string
			switch index {
			case 0:
				raw = f.PrimaryAccountNumber
			case 1:
				if f.ExpirationDate != nil {
					raw = f.ExpirationDate.Format(expiryDateFormat)
				}
			case 2:
				raw = f.ServiceCode
			case 3:
				raw = f.DiscretionaryData
			}

			if len(raw) == 0 {
				return `=`
			}
			return raw
		}
	case VersionThird:
		pattern = trackThirdPattern
		handler = func(index int, name string, group *syntax.Regexp, generator utils.Generator, args *utils.GeneratorArgs) string {
			var raw string
			switch index {
			case 0:
				raw = f.FormatCode
			case 1:
				raw = f.PrimaryAccountNumber
			case 2:
				raw = f.DiscretionaryData
			}

			if len(raw) == 0 {
				return `=`
			}
			return raw
		}
	}

	generator, _ := utils.NewGenerator(pattern, &utils.GeneratorArgs{
		Flags:               syntax.Perl,
		CaptureGroupHandler: handler,
	})

	rawTrack := generator.Generate()
	if matched, _ := regexp.MatchString(pattern, rawTrack); !matched {
		fmt.Println(rawTrack)
		return nil, errors.New("unable to create valid track data")
	}

	return []byte(rawTrack), nil
}
