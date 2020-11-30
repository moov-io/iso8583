// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package utils

import "regexp"

const (
	ElementTypeMti                 = "mti" // numeric characters
	ElementTypeBitmap              = "bit" // numeric characters
	ElementTypeAlphabetic          = "a"   // alphabetic characters only
	ElementTypeNumeric             = "n"   // numeric characters only
	ElementTypeSpecial             = "s"   // special characters only
	ElementTypeMagnetic            = "z"   // magnetic stripe track-2 or track-3 data
	ElementTypeIndicate            = "x"   // character “C” or “D” to indicate “credit” or “debit” value of a dollar amount
	ElementTypeBinary              = "b"   // binary data
	ElementTypeAlphaNumeric        = "an"  // alpha and numeric characters
	ElementTypeAlphaSpecial        = "as"  // alpha and special characters
	ElementTypeNumericSpecial      = "ns"  // numeric and special characters
	ElementTypeAlphaNumericSpecial = "ans" // alpha, numeric, and special characters
	ElementTypeIndicateNumeric     = "x+n" // Numeric (amount) values, where the first byte is either “C” or “D”
	ElementTypeNumberEncoding      = "number"

	DataElementXmlName    = "Element"
	DataElementAttrNumber = "Number"

	EncodingChar   = "CHAR"
	EncodingHex    = "HEX"
	EncodingEbcdic = "EBCDIC"
	EncodingAscii  = "ASCII"
	EncodingBcd    = "BCD" // packed bcd
	EncodingRBcd   = "RBCD"

	EncodingCatNumber    = "number"
	EncodingCatBinary    = "binary"
	EncodingCatCharacter = "character"
)

var (
	RegexAlphabetic          = regexp.MustCompile(`^[a-z A-Z]+$`).MatchString
	RegexNumeric             = regexp.MustCompile(`^[0-9]+$`).MatchString
	RegexSpecial             = regexp.MustCompile(`^[$&+,:;=?@#|'<>.^*()%! -]+$`).MatchString
	RegexIndicate            = regexp.MustCompile(`^[C|D]{1}$`).MatchString
	RegexAlphaNumeric        = regexp.MustCompile(`^[a-z A-Z0-9]*$`).MatchString
	RegexIndicateNumeric     = regexp.MustCompile(`^[C|D][0-9]+$`).MatchString
	RegexAlphaSpecial        = regexp.MustCompile(`^[a-zA-Z$&+,:;=?@#|'<>.^*()%! -]+$`).MatchString
	RegexBinary              = regexp.MustCompile(`^[0|1]+$`).MatchString
	RegexNumericSpecial      = regexp.MustCompile(`^[0-9$&+,:;=?@#|'<>.^*()%! -]+$`).MatchString
	RegexAlphaNumericSpecial = regexp.MustCompile(`^[0-9a-zA-Z$&+,:;=?@#|'<>.^*()%! -]+$`).MatchString

	RegexTimeHHMMSS = regexp.MustCompile(`(2[0-3]|[01][0-9])[0-5][0-9][0-5][0-9]`).MatchString
	RegexDateYYMM   = regexp.MustCompile(`((\d{2})(0[1-9]|10|12))`).MatchString
	RegexDateMMDD   = regexp.MustCompile(`((0[1-9]|10|12)(0[1-9]|[12][0-9]|3[01]))`).MatchString
	RegexDateYYMMDD = regexp.MustCompile(`((\d{2})(0[1-9]|10|11|12)(0[1-9]|[12][0-9]|3[01]))`).MatchString
)

// data representation attributes
var ElementDataTypes = []string{
	ElementTypeAlphabetic,
	ElementTypeNumeric,
	ElementTypeSpecial,
	ElementTypeMagnetic,
	ElementTypeIndicate,
	ElementTypeBinary,
	ElementTypeAlphaNumeric,
	ElementTypeAlphaSpecial,
	ElementTypeNumericSpecial,
	ElementTypeAlphaNumericSpecial,
	ElementTypeIndicateNumeric,
}

// data representation attributes
var AvailableEncodings = map[string][]string{
	ElementTypeMti:                 {EncodingChar, EncodingBcd},
	ElementTypeBitmap:              {EncodingChar, EncodingHex},
	ElementTypeAlphabetic:          {EncodingAscii, EncodingEbcdic},
	ElementTypeNumeric:             {EncodingBcd, EncodingRBcd, EncodingChar},
	ElementTypeSpecial:             {EncodingAscii, EncodingEbcdic},
	ElementTypeMagnetic:            {EncodingAscii, EncodingEbcdic},
	ElementTypeIndicate:            {EncodingAscii, EncodingEbcdic},
	ElementTypeBinary:              {EncodingChar, EncodingHex},
	ElementTypeAlphaNumeric:        {EncodingAscii, EncodingEbcdic},
	ElementTypeAlphaSpecial:        {EncodingAscii, EncodingEbcdic},
	ElementTypeNumericSpecial:      {EncodingAscii, EncodingEbcdic},
	ElementTypeAlphaNumericSpecial: {EncodingAscii, EncodingEbcdic},
	ElementTypeIndicateNumeric:     {EncodingAscii, EncodingEbcdic},
	ElementTypeNumberEncoding:      {EncodingChar, EncodingHex, EncodingRBcd, EncodingBcd},
}

var AvailableTypeCategory = map[string]string{
	ElementTypeMti:                 EncodingCatNumber,
	ElementTypeBitmap:              EncodingCatBinary,
	ElementTypeAlphabetic:          EncodingCatCharacter,
	ElementTypeNumeric:             EncodingCatNumber,
	ElementTypeSpecial:             EncodingCatCharacter,
	ElementTypeMagnetic:            EncodingCatCharacter,
	ElementTypeIndicate:            EncodingCatCharacter,
	ElementTypeBinary:              EncodingCatBinary,
	ElementTypeAlphaNumeric:        EncodingCatCharacter,
	ElementTypeAlphaSpecial:        EncodingCatCharacter,
	ElementTypeNumericSpecial:      EncodingCatCharacter,
	ElementTypeAlphaNumericSpecial: EncodingCatCharacter,
	ElementTypeIndicateNumeric:     EncodingCatCharacter,
}

var AvailableDateFormat = map[string]func(s string) bool{
	"HHMMSS": RegexTimeHHMMSS,
	"YYMM":   RegexDateYYMM,
	"MMDD":   RegexDateMMDD,
	"YYMMDD": RegexDateYYMMDD,
}

func CheckAvailableEncoding(eType string, encoding string) bool {
	encodings, exit := AvailableEncodings[eType]
	if !exit {
		return false
	}
	for _, _encoding := range encodings {
		if _encoding == encoding {
			return true
		}
	}
	return false
}
