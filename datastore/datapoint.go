package datastore

type Datapoint interface {
	Ider

	Label() string
	SetLabel(label string)
	Key() string
	SetKey(key string)

	ValueType() (ValueType, error)
}
