package forms

type TextBox struct {
	*baseField `schema:"-"`
	Value      string
}

func NewTextbox(label, value string) *TextBox {
	t := &TextBox{
		baseField: &baseField{
			Label: label,
		},
		Value: value,
	}
	return t
}

func (tb TextBox) ValidationChecks() []*ValidationCheck {
	return tb.validationRules
}

func (tb TextBox) SetError(message string) {
	tb.Error = message
}

func (tb TextBox) AddValidation(message string, check func(interface{}, interface{}) bool) {
	tb.validationRules = append(tb.validationRules, &ValidationCheck{message, check})
}
