package datastore

type DatastoreProvider interface {
	Initialize()

	// Subject Methods
	CreateSubject(subjectType, name, sortName string) Subject
	LoadSubject(id int) (Subject, error)

	// Datapoint Methods
	CreateDatapoint(label, key string, subjectTypeId, PrimativeTypeId int) Datapoint
	LoadDatapoint(id int) (Datapoint, error)

	// ValueType Methods
	CreateValueType(name string) ValueType
	LoadValueType(id int) (ValueType, error)

	// ValueType Field Methods
	CreateValueTypeField(label, key string, valueTypeId, primativeTypeId int) ValueTypeField
	//LoadValueTypeField(id int) (ValueTypeField error)

	// Primitive Types
	CreatePrimitiveType(name, label string) PrimitiveType
	LoadPrimitiveType(id int) (PrimitiveType, error)
}

var Provider DatastoreProvider

func Register(provider DatastoreProvider) {
	provider.Initialize()
	Provider = provider
}
