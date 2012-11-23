package datastore

// Defines an interface for thing that should have a unique ID
type Ider interface {
	Id() int
}

// An interface for deleting datastore objects
type Deleter interface {
	Delete() error
}

// An interface for saving datastore objects
type Saver interface {
	Save() error
}
