package forms

type SingleCheckbox struct {
	*baseField `schema:"-"`
	Checked    bool
}

func NewCheckbox(label string, checked bool) *SingleCheckbox {
	return &SingleCheckbox{
		baseField: &baseField{
			Label: label,
		},
		Checked: checked,
	}
}

func (field SingleCheckbox) ValidationChecks() []*ValidationCheck {
	return field.validationRules
}

func (field SingleCheckbox) SetError(message string) {
	field.Error = message
}

func (field SingleCheckbox) AddValidation(message string, check func(interface{}, interface{}) bool) {
	field.validationRules = append(field.validationRules, &ValidationCheck{message, check})
}
