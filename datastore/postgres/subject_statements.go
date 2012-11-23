package postgres

import (
	"database/sql"
)

type subjectStatements struct {
	load   *sql.Stmt
	update *sql.Stmt
	create *sql.Stmt
	delete *sql.Stmt

	updateData *sql.Stmt
	deleteData *sql.Stmt
}

func newSubjectStatements(pg *PostgresProvider) subjectStatements {

	s := subjectStatements{}

	s.load = pg.mustPrepare("SELECT id, type, name, sort_name, data FROM subject WHERE id = $1")
	s.update = pg.mustPrepare("UPDATE subject SET name = $2, sort_name = $3 WHERE id = $1")
	s.create = pg.mustPrepare("INSERT INTO subject (type, name, sort_name, data) VALUES ($1, $2, $3, $4) RETURNING id")
	s.delete = pg.mustPrepare("DELETE FROM subject WHERE id = $1")

	s.updateData = pg.mustPrepare("UPDATE subject SET data = data || hstore($1, $2) WHERE id = $3")
	s.deleteData = pg.mustPrepare("UPDATE subject SET data = delete(data, $1) WHERE id = $2")

	return s
}
