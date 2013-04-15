package forms

import (
	"errors"
	"fmt"
	"reflect"
)

func ValidatePresence(form interface{}, field interface{}) bool {
	switch t := field.(type) {
	case *TextBox:
		return t.Value != ""
	case *Radios:
		return t.Value != ""
	case *Textarea:
		return t.Value != ""
	default:
		fmt.Printf("Cannot validate presence of %T", t)
	}
	return true
}

type ValidationError error

func Validate(form interface{}) []ValidationError {
	errs := make([]ValidationError, 0)

	formValue := reflect.ValueOf(form).Elem()

	for i := 0; i < formValue.NumField(); i++ {
		field := formValue.Field(i)
		if field.CanInterface() {
			if v, ok := field.Interface().(Validatable); ok {
				for _, check := range v.ValidationChecks() {
					if !check.check(form, field.Interface()) {
						fmt.Println("checking ", field.Interface())
						v.SetError(check.message)
						errs = append(errs, errors.New(check.message))
					}
				}
			}
		}
	}

	return errs
}
