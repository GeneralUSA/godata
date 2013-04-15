package forms

type Checkboxes struct {
	*baseField
	Options []*Option
}

func NewCheckboxes(label string) *Checkboxes {
	c := &Checkboxes{
		baseField: &baseField{
			Label: label,
		},
		Options: make([]*Option, 0),
	}
	return c
}

func (cb *Checkboxes) AddOption(value, label string, checked bool) {
	cb.Options = append(cb.Options, &Option{value, label, false, checked})
}

func (cb Checkboxes) AddValidation(message string, check func(interface{}, interface{}) bool) {
	cb.validationRules = append(cb.validationRules, &ValidationCheck{message, check})
}

func (cb Checkboxes) ValidationChecks() []*ValidationCheck {
	return cb.validationRules
}

func (cb Checkboxes) SetError(err string) {
	cb.Error = err
}
