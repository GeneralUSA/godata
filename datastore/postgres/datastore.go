package postgres

import (
	"code.google.com/p/goconf/conf"
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
)

var db *sql.DB

func getConnection() *sql.DB {
	if db == nil {
		var err error
		c, err := conf.ReadConfigFile("../database.config")
		if err != nil {
			panic(err)
		}

		host, _ := c.GetString("default", "host")
		user, _ := c.GetString("default", "user")
		pass, _ := c.GetString("default", "password")
		dbname, _ := c.GetString("default", "database")

		connectionString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v", host, user, pass, dbname)

		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
	}
	return db
}
