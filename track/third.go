package track

import (
	"errors"
	"regexp"
	"regexp/syntax"
	"strings"
)

var _ Track = (*Third)(nil)

type Third struct{}

func NewTrackThird() *Third {
	return &Third{}
}

func (h *Third) Write(card *GeneralCard) ([]byte, error) {
	generator, _ := NewGenerator(trackThirdPattern, &GeneratorArgs{
		Flags: syntax.Perl,
		CaptureGroupHandler: func(index int, name string, group *syntax.Regexp, generator Generator, args *GeneratorArgs) string {
			var raw string
			switch index {
			case 0:
				raw = card.FormatCode
			case 1:
				raw = card.PrimaryAccountNumber
			case 2:
				raw = card.DiscretionaryData
			}

			if len(raw) == 0 {
				return `=`
			}
			return raw
		},
	})

	rawTrack := generator.Generate()
	if matched, _ := regexp.MatchString(trackThirdPattern, rawTrack); !matched {
		return nil, errors.New("unable to create valid track data")
	}

	if len(rawTrack) > trackThirdMaxLength {
		return nil, errors.New("unable to create valid track data")
	}

	return []byte(rawTrack), nil
}

func (h *Third) Read(raw []byte) (*GeneralCard, error) {
	if raw == nil || len(raw) > trackThirdMaxLength {
		return nil, errors.New("invalid track 3 format")
	}

	r, err := regexp.Compile(trackThirdPattern)
	if err != nil {
		return nil, err
	}

	if !r.MatchString(string(raw)) {
		return nil, errors.New("invalid track 3 format")
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
			card.FormatCode = value
		case 2: // Payment card number (PAN)
			card.PrimaryAccountNumber = value
			card.CardType = GetCardType(value)
		case 3: // Security Data + Additional Data
			card.DiscretionaryData = value
		}
	}

	return &card, nil
}
