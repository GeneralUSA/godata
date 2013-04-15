package types

import (
	"fmt"
	"net/url"
	"reflect"
)

type singleValue struct {
	name        string
	label       string
	description string
	validators  []*Validator
	template    string
	fieldType   FieldType
}

func (sv *singleValue) Name() string {
	return sv.name
}

func (sv *singleValue) Label() string {
	return sv.label
}

func (sv *singleValue) Description() string {
	return sv.description
}

func (sv *singleValue) Template() string {
	return sv.template
}

func (sv *singleValue) FieldType() FieldType {
	return sv.fieldType
}

func (sv *singleValue) Parse(form url.Values) (reflect.Value, error) {
	if values, ok := form["Value"]; ok && len(values) > 0 {
		return sv.fieldType.Parse(values[0])
	}
	return reflect.Value{}, fmt.Errorf("No value")
}
