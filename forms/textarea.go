package forms

type Textarea struct {
	*baseField `schema:"-"`
	Value      string
	Rows       int `schema:"-"`
}

func NewTextarea(label, value string) *Textarea {
	return &Textarea{
		baseField: &baseField{
			Label: label,
		},
		Value: value,
		Rows:  2,
	}
}

func (field Textarea) ValidationChecks() []*ValidationCheck {
	return field.validationRules
}

func (field Textarea) SetError(message string) {
	field.Error = message
}

func (field Textarea) AddValidation(message string, check func(interface{}, interface{}) bool) {
	field.validationRules = append(field.validationRules, &ValidationCheck{message, check})
}
