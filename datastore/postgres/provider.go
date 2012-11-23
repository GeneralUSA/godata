package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
)

type PostgresProvider struct {
	host       string
	user       string
	password   string
	dbname     string
	db         *sql.DB
	statements struct {
		subject        subjectStatements
		primitiveType  primativeTypeStatements
		valueType      valueTypeStatements
		valueTypeField valueTypeFieldStatements
		datapoint      datapointStatements
	}
}

// Creates a new Datastore provider for using a Postgres Database
// as the backend.
func NewPostgresProvidier(host, user, password, dbname string) *PostgresProvider {
	return &PostgresProvider{
		host:     host,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

var db *sql.DB
var provider *PostgresProvider

func (pg *PostgresProvider) Initialize() {

	connectionString := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v",
		pg.host,
		pg.user,
		pg.password,
		pg.dbname,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	pg.db = db
	pg.statements.subject = newSubjectStatements(pg)
	pg.statements.primitiveType = newPrimativeTypeStatements(pg)
	pg.statements.valueType = newValueTypeStatements(pg)
	pg.statements.valueTypeField = newValueTypeFieldStatements(pg)
	pg.statements.datapoint = newDatapointStatements(pg)
	provider = pg
}

// Attempts to prepare a statement on the provider's database connection,
// if the statement fails to prepare, the function panics.
func (pg PostgresProvider) mustPrepare(query string) *sql.Stmt {
	stmt, err := pg.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}
