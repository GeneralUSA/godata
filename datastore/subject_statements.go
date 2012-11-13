package datastore

import (
	"database/sql"
)

type SubjectStatements struct {
	load   *sql.Stmt
	update *sql.Stmt
	create *sql.Stmt
	delete *sql.Stmt

	updateData *sql.Stmt
	deleteData *sql.Stmt
}

func newSubjectStatements(db *sql.DB) *SubjectStatements {
	s := SubjectStatements{}
	s.load, _ = db.Prepare("SELECT id, type, name, sort_name, data FROM subject WHERE id = $1")
	s.update, _ = db.Prepare("UPDATE subject SET name = $2, sort_name = $3 WHERE id = $1")
	s.create, _ = db.Prepare("INSERT INTO subject (type, name, sort_name, data) VALUES ($1, $2, $3, $4) RETURNING id")
	s.delete, _ = db.Prepare("DELETE FROM subject WHERE id = $1")

	s.updateData, _ = db.Prepare("UPDATE subject SET data = data || hstore($1, $2) WHERE id = $3")
	s.deleteData, _ = db.Prepare("UPDATE subject SET data = delete(data, $1) WHERE id = $2")

	return &s
}
