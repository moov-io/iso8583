package track

import (
	"errors"
	"regexp"
	"regexp/syntax"
	"strings"
)

var _ Track = (*Second)(nil)

type Second struct{}

func NewTrackSecond() *Second {
	return &Second{}
}

func (h *Second) Write(card *GeneralCard) ([]byte, error) {
	generator, _ := NewGenerator(trackSecondPattern, &GeneratorArgs{
		Flags: syntax.Perl,
		CaptureGroupHandler: func(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string {
			var raw string
			switch index {
			case 0:
				raw = card.PrimaryAccountNumber
			case 1:
				if card.ExpirationDate != nil {
					raw = card.ExpirationDate.String()
				}
			case 2:
				raw = card.ServiceCode
			case 3:
				raw = card.DiscretionaryData
			}

			if len(raw) == 0 {
				return `=`
			}
			return raw
		},
	})

	rawTrack := generator.Generate()
	if matched, _ := regexp.MatchString(trackSecondPattern, rawTrack); !matched {
		return nil, errors.New("unable to create valid track data")
	}

	if len(rawTrack) > trackSecondMaxLength {
		return nil, errors.New("unable to create valid track data")
	}

	return []byte(rawTrack), nil
}

func (h *Second) Read(raw []byte) (*GeneralCard, error) {
	if raw == nil || len(raw) > trackSecondMaxLength {
		return nil, errors.New("invalid track 2 format")
	}

	r, err := regexp.Compile(trackSecondPattern)
	if err != nil {
		return nil, err
	}

	if !r.MatchString(string(raw)) {
		return nil, errors.New("invalid track 2 format")
	}

	var card GeneralCard
	matches := r.FindStringSubmatch(string(raw))
	for index, val := range matches {
		value := strings.TrimSpace(val)
		if len(value) == 0 || value == "=" {
			continue
		}

		switch index {
		case 1: // Payment card number (PAN)
			card.PrimaryAccountNumber = value
			card.CardType = GetCardType(value)
		case 2: // Expiration Date (ED)
			card.ExpirationDate, err = NewExpiryDate(value)
		case 3: // Service Code (SC)
			card.ServiceCode = value
		case 4: // Discretionary data (DD)
			card.DiscretionaryData = value
		}
	}

	return &card, err
}
