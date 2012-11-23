package postgres

import (
	"database/sql"
	"fmt"
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type ValueType struct {
	id     int
	name   string
	fields []datastore.ValueTypeField

	isNew bool
}

type valueTypeStatements struct {
	create *sql.Stmt
	update *sql.Stmt
	load   *sql.Stmt
	delete *sql.Stmt
}

func newValueTypeStatements(pg *PostgresProvider) valueTypeStatements {
	s := valueTypeStatements{}

	s.create = pg.mustPrepare("INSERT INTO value_type (name) VALUES ($1) RETURNING id")
	s.update = pg.mustPrepare("UPDATE value_type SET name = $2 WHERE id = $1")
	s.load = pg.mustPrepare("SELECT id, name FROM value_type WHERE id = $1")
	s.delete = pg.mustPrepare("DELETE FROM value_type WHERE id = $1")

	return s
}

func (pg PostgresProvider) CreateValueType(name string) datastore.ValueType {
	vt := ValueType{
		name:  name,
		isNew: true,
	}
	return &vt
}

func (pg PostgresProvider) LoadValueType(id int) (datastore.ValueType, error) {
	var err error
	row := pg.statements.valueType.load.QueryRow(id)
	vt := ValueType{}
	err = row.Scan(&vt.id, &vt.name)
	if err != nil {
		return nil, err
	}

	vt.fields, err = pg.loadAllForValueType(vt.id)
	if err != nil {
		return nil, err
	}
	return &vt, nil
}

func (v *ValueType) Save() error {
	if v.isNew {
		return v.insert()
	} else {
		return v.update()
	}
	return nil
}

func (v *ValueType) insert() error {
	row := provider.statements.valueType.create.QueryRow(v.name)
	err := row.Scan(&v.id)
	if err != nil {
		return err
	}
	v.isNew = false
	return nil
}

func (v ValueType) update() error {
	_, err := provider.statements.valueType.update.Exec(v.id, v.name)
	return err
}

func (v ValueType) Delete() error {
	_, err := provider.statements.valueType.delete.Exec(v.id)
	return err
}

func (v ValueType) Fields() []datastore.ValueTypeField {
	return v.fields
}

func (v ValueType) Id() int {
	return v.id
}

func (v ValueType) Name() string {
	return v.name
}

func (v ValueType) SetName(name string) {
	v.name = name
}

func (v ValueType) String() string {
	s := v.name
	for t := range v.fields {
		s += fmt.Sprintf("\n  %v", v.fields[t])
	}
	return s
}
