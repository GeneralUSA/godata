package datastore

type Subject interface {
	// Fields set by the datastore
	Id() int

	// Immutable Fields
	SubjectType() string

	// Mutable Fields
	Name() string
	SetName(name string)
	SortName() string
	SetSortName(sortName string)

	// Operations
	Save() error
	Delete() error
}
