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
		subject SubjectStatements
	}
}

func NewPostgresProvidier(host, user, password, dbname string) PostgresProvider {
	return PostgresProvider{
		host:     host,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

var db *sql.DB
var provider PostgresProvider

func (pg PostgresProvider) Initialize() {

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
	pg.statements.subject = newSubjectStatements(pg.db)
	provider = pg
}
