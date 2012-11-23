package datastore

type ValueType interface {
	Ider

	Name() string
	SetName(name string)

	Fields() []ValueTypeField

	Saver
	Deleter
}
