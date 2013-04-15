// Package global contains system-wide variables
package global

import (
	"database/sql"
	"github.com/coopernurse/gorp"
)

var (
	DB    *sql.DB
	DbMap *gorp.DbMap
)
