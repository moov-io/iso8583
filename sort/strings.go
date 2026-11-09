package sort

import (
	"sort"
	"strconv"
	"strings"
)

// StringSlice is a function type used to sort a slice of strings in increasing
// order. Any errors which arise from sorting the slice will raise a panic.
type StringSlice func(x []string)

// Strings sorts a slice of strings in increasing order.
var Strings = sort.Strings

// StringsByInt sorts a slice of strings according to their integer value.
// This function panics in the event that an element in the slice cannot be
// converted to an integer
func StringsByInt(x []string) {
	sort.Slice(x, func(i, j int) bool {
		valI, err := strconv.Atoi(x[i])
		if err != nil {
			return x[i] < x[j]
		}
		valJ, err := strconv.Atoi(x[j])
		if err != nil {
			return x[i] < x[j]
		}
		return valI < valJ
	})
}

// StringsByHex sorts a slice of strings according to their big-endian Hex
// value. This function compares values as a regular strings in the event that
// an element in the slice cannot be converted to a Hex slice. Each string
// representation of a hex value must be of even length.
func StringsByHex(x []string) {
	sort.Slice(x, func(i, j int) bool {
		return compareHex(x[i], x[j]) < 0
	})
}

// compareHex compares a and b as big-endian hexadecimal numbers. If either
// value is not valid hexadecimal, compareHex compares the strings
// lexicographically.
func compareHex(a, b string) int {
	aValue, aOK := trimHex(a)
	bValue, bOK := trimHex(b)
	if !aOK || !bOK {
		return strings.Compare(a, b)
	}

	if len(aValue) != len(bValue) {
		return len(aValue) - len(bValue)
	}

	for i := 0; i < len(aValue); i++ {
		aNibble, _ := hexNibble(aValue[i])
		bNibble, _ := hexNibble(bValue[i])
		if aNibble != bNibble {
			return int(aNibble) - int(bNibble)
		}
	}

	return 0
}

// trimHex validates an even-length hexadecimal string and removes leading zero
// digits so that remaining values can be compared by length and digits.
func trimHex(s string) (string, bool) {
	if len(s)%2 != 0 {
		return s, false
	}

	firstNonZero := -1
	for i := 0; i < len(s); i++ {
		if _, ok := hexNibble(s[i]); !ok {
			return s, false
		}
		if firstNonZero == -1 && s[i] != '0' {
			firstNonZero = i
		}
	}

	if firstNonZero == -1 {
		return "", true
	}

	return s[firstNonZero:], true
}

func hexNibble(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}
