package utils

import "strconv"

type NumericSort []string

func (a NumericSort) Len() int      { return len(a) }
func (a NumericSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a NumericSort) Less(i, j int) bool {
	numLeft, errLeft := strconv.ParseUint(a[i], 10, 0)
	numRight, errRight := strconv.ParseUint(a[j], 10, 0)
	if errLeft == nil && errRight == nil {
		return numLeft < numRight
	}
	if errLeft == nil {
		return true
	}
	if errRight == nil {
		return false
	}
	return a[i] < a[j]
}
