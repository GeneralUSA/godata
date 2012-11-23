package postgres

import (
	"database/sql"
	"fmt"
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type ValueTypeField struct {
	id              int
	name            string
	key             string
	primitiveTypeId int
	valueTypeId     int

	isNew bool
}

type valueTypeFieldStatements struct {
	create *sql.Stmt
	delete *sql.Stmt

	allValueType *sql.Stmt
}

func newValueTypeFieldStatements(pg *PostgresProvider) valueTypeFieldStatements {
	s := valueTypeFieldStatements{}

	s.create = pg.mustPrepare("INSERT INTO value_type_field (value_type_id, primative_type_id, name, key) VALUES ($1, $2, $3, $4) RETURNING id")
	s.delete = pg.mustPrepare("DELETE FROM value_type_field WHERE id = $1")

	s.allValueType = pg.mustPrepare("SELECT id, value_type_id, primative_type_id, name, key FROM value_type_field WHERE value_type_id = $1")

	return s
}

func (p PostgresProvider) CreateValueTypeField(name, key string, primitiveTypeId, valueTypeId int) datastore.ValueTypeField {
	vtf := ValueTypeField{
		name:            name,
		key:             key,
		primitiveTypeId: primitiveTypeId,
		valueTypeId:     valueTypeId,
		isNew:           true,
	}
	return &vtf
}

func (p PostgresProvider) loadAllForValueType(valueTypeId int) ([]datastore.ValueTypeField, error) {
	rows, err := p.statements.valueTypeField.allValueType.Query(valueTypeId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	results := make([]datastore.ValueTypeField, 0, 1)

	var vtf *ValueTypeField
	for rows.Next() {
		vtf = &ValueTypeField{}
		rows.Scan(&vtf.id, &vtf.valueTypeId, &vtf.primitiveTypeId, &vtf.name, &vtf.key)
		results = append(results, vtf)
	}

	return results, nil
}

func (vtf ValueTypeField) Delete() error {
	_, err := provider.statements.valueTypeField.delete.Exec(vtf.id)
	return err
}

func (vtf *ValueTypeField) Save() error {
	if vtf.isNew {
		return vtf.insert()
	} else {
		return vtf.update()
	}
	return nil
}

func (vtf *ValueTypeField) insert() error {
	row := provider.statements.valueTypeField.create.QueryRow(vtf.valueTypeId, vtf.primitiveTypeId, vtf.name, vtf.key)
	err := row.Scan(&vtf.id)
	if err != nil {
		return err
	}
	vtf.isNew = false
	return nil
}

func (vtf *ValueTypeField) update() error {
	return nil
}

func (vtf ValueTypeField) Id() int {
	return vtf.id
}

func (vtf ValueTypeField) Name() string {
	return vtf.name
}

func (vtf *ValueTypeField) SetName(name string) {
	vtf.name = name
}

func (vtf ValueTypeField) Key() string {
	return vtf.key
}

func (vtf *ValueTypeField) SetKey(key string) {
	vtf.key = key
}

func (vtf *ValueTypeField) SetPrimitiveTypeId(id int) {
	vtf.primitiveTypeId = id
}

func (vtf *ValueTypeField) SetValueTypeId(id int) {
	vtf.valueTypeId = id
}

func (vtf ValueTypeField) String() string {
	ptype := vtf.PrimitiveType()
	return fmt.Sprintf("%v (%v)", vtf.name, ptype.Label())
}

func (vtf ValueTypeField) PrimitiveType() datastore.PrimitiveType {
	// BUG(andrew) This should not hit the database everytime its called
	t, _ := provider.LoadPrimitiveType(vtf.primitiveTypeId)
	return t
}
