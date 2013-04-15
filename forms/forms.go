package forms

import (
	"fmt"
	"github.com/gorilla/schema"
	"reflect"
	"strconv"
)

var decoder = schema.NewDecoder()

func DecodeForm(dest interface{}, src map[string][]string) error {
	return decoder.Decode(dest, src)
}

func convertBool(value string) reflect.Value {
	if value == "on" {
		return reflect.ValueOf(true)
	} else if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.ValueOf(false)
}

func init() {
	decoder.RegisterConverter(false, convertBool)
}

type Form interface {
}

type baseField struct {
	Id              string
	Name            string
	Label           string
	Value           string
	Help            string
	validationRules []*ValidationCheck
	Error           string
}

type ValidationCheck struct {
	message string
	check   func(interface{}, interface{}) bool
}

type Validatable interface {
	AddValidation(message string, check func(interface{}, interface{}) bool)
	ValidationChecks() []*ValidationCheck
	SetError(message string)
}

type Option struct {
	Key     string
	Label   string
	Checked bool
	Set     bool
}

func (o Option) String() string {
	return fmt.Sprintf("%v:%v (%v)", o.Key, o.Label, o.Checked)
}

func PrepareForm(form interface{}) {

	formValue := reflect.ValueOf(form).Elem()
	formType := reflect.TypeOf(form).Elem()

	for i := 0; i < formValue.NumField(); i++ {
		field := formValue.Field(i)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Kind() == reflect.Struct {
			nameField := field.FieldByName("Name")
			if nameField.CanSet() {
				fieldInfo := formType.Field(i)
				tag := fieldInfo.Tag.Get("schema")
				if tag != "" {
					nameField.SetString(tag)
				} else {
					nameField.SetString(fieldInfo.Name)
				}
			}
		}
	}
}

func NewOptions(options map[string]string) []*Option {
	o := make([]*Option, 0, len(options))
	for k, v := range options {
		o = append(o, &Option{Key: k, Label: v})
	}
	return o
}
