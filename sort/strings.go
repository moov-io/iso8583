package sort

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/moov-io/iso8583/encoding"
)

// StringSlice is a function type used to sort a slice of strings in increasing
// order. Any errors which arise from sorting the slice will raise a panic.
type StringSlice func(x []string) error

// Strings sorts a slice of strings in increasing order.
func Strings(x []string) error {
	sort.Strings(x)
	return nil
}

// StringsByInt sorts a slice of strings according to their integer value.
// This function panics in the event that an element in the slice cannot be
// converted to an integer
func StringsByInt(x []string) error {
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
	return nil
}

// StringsByHex sorts a slice of strings according to their big-endian Hex value.
// This function panics in the event that an element in the slice cannot be
// converted to a Hex slice. Each string representation of a hex value must be
// of even length.
func StringsByHex(x []string) error {
	var outerErr error
	sort.Slice(x, func(i, j int) bool {
		valI, err := encoding.ASCIIHexToBytes.Encode([]byte(x[i]))
		if err != nil {
			outerErr = fmt.Errorf("failed to encode ascii hex %s to bytes : %v", x[i], err)
			return false
		}
		valJ, err := encoding.ASCIIHexToBytes.Encode([]byte(x[j]))
		if err != nil {
			outerErr = fmt.Errorf("failed to sort strings by hex: %v", err)
			return false
		}
		return new(big.Int).SetBytes(valI).Int64() < new(big.Int).SetBytes(valJ).Int64()
	})
	return outerErr
}
