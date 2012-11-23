package postgres

import (
	"database/sql"
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type PrimitiveType struct {
	id    int
	name  string
	label string

	isNew bool
}

type primativeTypeStatements struct {
	load   *sql.Stmt
	create *sql.Stmt
	delete *sql.Stmt
	update *sql.Stmt
}

func newPrimativeTypeStatements(pg *PostgresProvider) primativeTypeStatements {
	s := primativeTypeStatements{}

	s.load = pg.mustPrepare("SELECT id, name, label from primative_type WHERE id = $1")
	s.create = pg.mustPrepare("INSERT INTO primative_type (name, label) VALUES ($1, $2) RETURNING id")
	s.update = pg.mustPrepare("UPDATE primative_type SET name = $2, label = $3 WHERE id = $1")
	s.delete = pg.mustPrepare("DELETE FROM primative_type WHERE id = $1")

	return s
}

func (pg PostgresProvider) LoadPrimitiveType(id int) (datastore.PrimitiveType, error) {
	pt := PrimitiveType{}
	row := pg.statements.primitiveType.load.QueryRow(id)
	err := row.Scan(&pt.id, &pt.name, &pt.label)
	if err != nil {
		return nil, err
	}
	return &pt, nil
}

func (pg PostgresProvider) CreatePrimitiveType(name, label string) datastore.PrimitiveType {
	pt := PrimitiveType{
		name:  name,
		label: label,
		isNew: true,
	}
	return &pt
}

func (p PrimitiveType) Id() int {
	return p.id
}

func (p PrimitiveType) Label() string {
	return p.label
}

func (p *PrimitiveType) SetLabel(label string) {
	p.label = label
}

func (p PrimitiveType) Name() string {
	return p.name
}

func (p *PrimitiveType) SetName(name string) {
	p.name = name
}

func (p *PrimitiveType) Save() error {
	if p.isNew {
		return p.insert()
	} else {
		return p.update()
	}
	return nil
}

func (p PrimitiveType) update() error {
	_, err := provider.statements.primitiveType.update.Exec(p.id, p.name, p.label)
	return err
}

func (p *PrimitiveType) insert() error {
	row := provider.statements.primitiveType.create.QueryRow(p.name, p.label)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	p.id = id
	p.isNew = false
	return nil
}

func (p PrimitiveType) Delete() error {
	_, err := provider.statements.primitiveType.delete.Exec(p.id)
	return err
}

func (p PrimitiveType) String() string {
	return p.name
}
