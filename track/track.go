package track

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

/*
Track is the interface for card track data

Track provides utility functions for reading and writing magnetic track data
http://sagan.gae.ucm.es/~padilla/extrawork/tracks.html
https://en.wikipedia.org/wiki/ISO/IEC_7813
https://www.q-card.com/about-us/iso-magnetic-stripe-card-standards/page.aspx?id=1457


UseCase

message := iso8583.NewMessage(spec)
message.MTI("0100")
message.Field(2, "4242424242424242")
message.Field(3, "123456")
message.Field(4, "100")
message.Field(4, "100")

dataBuf, _ := NewTrackSecond().Write(card)
message.Field(36, string(dataBuf))

*/

type GeneralCard struct {
	// Available for Track 1, Track 3
	FormatCode string `json:"format_code"`

	// Available for Track 1, Track 2, Track 3
	PrimaryAccountNumber string `json:"primary_account_number"`

	// Available for Track 1, Track 2, Track 3
	CardType string `json:"card_type"`

	// Available for Track 1
	Name string `json:"name"`

	// Available for Track 1, Track 2
	ExpirationDate *ExpiryDate `json:"expiration_date"`

	// Available for Track 1, Track 2
	ServiceCode string `json:"expiration_date"`

	// Available for Track 1, Track 2, Track 3
	//  If track 3, the field describe security data + additional data
	DiscretionaryData string `json:"discretionary_data"`
}

// Track is network header interface to write/read encoded message legnth
type Track interface {
	// Write general card into raw track date
	Write(*GeneralCard) ([]byte, error)

	// Read reads general card from raw track data
	Read([]byte) (*GeneralCard, error)
}

type CardType struct {
	Pattern string
	Name    string
}

var CardTypes = map[string]string{
	"Visa":             `^4[0-9]{12}(?:[0-9]{3})?$`,
	"MasterCard":       `^5[1-5][0-9]{14}$`,
	"American Express": `^3[47][0-9]{13}$`,
	"Diners Club":      `^3(?:0[0-5]|[68][0-9])[0-9]{11}$`,
	"Discover":         `^6(?:011|5[0-9]{2})[0-9]{12}$`,
	"JCB":              `^(?:2131|1800|35\d{3})\d{11}$`,
}

func GetCardType(number string) string {
	for key, pattern := range CardTypes {
		r, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		if !r.MatchString(number) {
			return key
		}
	}
	return "Unknown"
}

// track pattern doesn't contain SS, ES, LRC
const (
	expiryDatePattern    = `^([0-9]{2})\/?(0[1-9]|1[0-2])$` //YYMM
	trackFirstPattern    = `^([A-Z]{1})([0-9]{1,19})\^([^\^]{2,26})\^([0-9]{4}|\^)([0-9]{3}|\^)([^\?]+)$`
	trackSecondPattern   = `^([0-9]{1,19})\=([0-9]{4}|\=)([0-9]{3}|\=)([^\?]+)$`
	trackThirdPattern    = `^([0-9]{2})([0-9]{1,19})\=([^\?]+)$`
	trackFirstMaxLength  = 76
	trackSecondMaxLength = 37
	trackThirdMaxLength  = 104
)

type ExpiryDate time.Time

func NewExpiryDate(str string) (*ExpiryDate, error) {
	if len(str) == 0 {
		return nil, nil //skip
	}

	if len(str) == 0 {
		return nil, errors.New("invalid expired date string")
	}

	str = str[:2] + "/" + str[2:4]
	r, err := regexp.Compile(expiryDatePattern)
	if err != nil {
		return nil, errors.New("unable to match date pattern")
	}

	matches := r.FindStringSubmatch(str)
	if len(matches) == 0 {
		return nil, errors.New("unable to match date pattern")
	}
	month, _ := strconv.Atoi(matches[2])
	year, _ := strconv.Atoi(matches[1])
	expired := ExpiryDate(time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC))

	return &expired, nil
}

func (t ExpiryDate) String() string {
	if (time.Time)(t).IsZero() {
		return ""
	}
	return fmt.Sprintf("%2s%02d", strconv.Itoa((time.Time)(t).Year())[2:4], (time.Time)(t).Month())
}
