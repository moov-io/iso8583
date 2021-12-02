package sort

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/moov-io/iso8583/encoding"
)

// Strings is a function type used to sort a slice of strings in increasing
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
			panic("failed to sort strings by int: failed to convert string to int")
		}
		valJ, err := strconv.Atoi(x[j])
		if err != nil {
			panic("failed to sort strings by int: failed to convert string to int")
		}
		return valI < valJ
	})
}

// StringsByHex sorts a slice of strings according to their big-endian Hex value.
// This function panics in the event that an element in the slice cannot be
// converted to a Hex slice. Each string representation of a hex value must be
// of even length.
func StringsByHex(x []string) {
	sort.Slice(x, func(i, j int) bool {
		valI, err := encoding.ASCIIHexToBytes.Encode([]byte(x[i]))
		if err != nil {
			panic(fmt.Sprintf("failed to sort strings by hex: %v", err))
		}
		valJ, err := encoding.ASCIIHexToBytes.Encode([]byte(x[j]))
		if err != nil {
			panic(fmt.Sprintf("failed to sort strings by hex: %v", err))
		}
		return new(big.Int).SetBytes(valI).Int64() < new(big.Int).SetBytes(valJ).Int64()
	})
}
