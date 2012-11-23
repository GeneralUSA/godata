package datastore

type PrimitiveType interface {
	Ider

	Label() string
	SetLabel(label string)

	Name() string
	SetName(name string)

	Saver
	Deleter
}
