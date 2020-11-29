package validation

import (
	"fmt"
	"strconv"
)

type Validator interface {
	IsValid([]byte) error
}

var None Validator = &noneValidator{}

type noneValidator struct{}

func (*noneValidator) IsValid(data []byte) error {
	return nil
}

var Numeric Validator = &numericValidator{}

type numericValidator struct{}

func (*numericValidator) IsValid(data []byte) error {
	a := string(data)

	if _, err := strconv.Atoi(a); err != nil {
		return fmt.Errorf("Expected: %s to be a number", a)
	}

	return nil
}
