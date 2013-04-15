package forms

type Radios struct {
	*baseField
	Value   string
	Options []*Option
}

func NewRadios(label string) *Radios {
	r := &Radios{
		baseField: &baseField{
			Label: label,
		},
		Options: make([]*Option, 0),
	}

	return r
}

func (r *Radios) AddOption(value, label string, checked bool) {
	r.Options = append(r.Options, &Option{value, label, checked, checked})
}

func (r Radios) AddValidation(message string, check func(interface{}, interface{}) bool) {
	r.validationRules = append(r.validationRules, &ValidationCheck{message, check})
}

func (r Radios) ValidationChecks() []*ValidationCheck {
	return r.validationRules
}

func (r Radios) SetError(err string) {
	r.Error = err
}
