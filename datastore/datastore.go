package datastore

type DatastoreProvider interface {
	Initialize()

	// Subject Methods
	CreateSubject(subjectType, name, sortName string) Subject
	LoadSubject(id int) (Subject, error)
}

var Provider DatastoreProvider

func Register(provider DatastoreProvider) {
	provider.Initialize()
	Provider = provider
}
