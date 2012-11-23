package datastore

type Subject interface {
	// Fields set by the datastore
	Ider

	// Immutable Fields
	SubjectType() string

	// Mutable Fields
	Name() string
	SetName(name string)
	SortName() string
	SetSortName(sortName string)

	// Operations
	Saver
	Deleter
}
