// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package pkg

const (
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
	DataElementXmlName             = "Element"
	DataElementAttrNumber          = "Number"
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
