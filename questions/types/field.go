package types

import (
	"fmt"
	"reflect"
	"strconv"
)

type FieldType int

const (
	FieldString = FieldType(iota)
	FieldInteger
)

func (t FieldType) Table() string {
	switch t {
	case FieldString:
		return "data_string"
	case FieldInteger:
		return "data_int"
	}
	return ""
}

func (t FieldType) Parse(value string) (reflect.Value, error) {
	switch t {
	case FieldString:
		return reflect.ValueOf(value), nil
	case FieldInteger:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int(i)), nil
	}
	return reflect.Value{}, fmt.Errorf("Unsupported Field Type")
}
