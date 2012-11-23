package datastore

type ValueTypeField interface {
	Ider

	Key() string
	SetKey(key string)

	SetValueTypeId(valueTypeId int)

	PrimitiveType() PrimitiveType
	SetPrimitiveTypeId(primitiveTypeId int)

	Saver
	Deleter
}
