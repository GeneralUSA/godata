package types

import (
	"strconv"
)

type Validator struct {
	Name     string
	Label    string
	Validate func(string) bool
}

var ValidateInt = &Validator{
	Name:     "integer",
	Label:    "Integer",
	Validate: validateInt,
}

func validateInt(value string) bool {
	_, err := strconv.ParseInt(value, 10, 32)
	return err == nil
}

var ValidateFloat = &Validator{
	Name:     "float",
	Label:    "Floating Point",
	Validate: validateFloat,
}

func validateFloat(value string) bool {
	_, err := strconv.ParseFloat(value, 32)
	return err == nil
}
