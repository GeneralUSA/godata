package users

import (
	"database/sql"
	"github.com/gorilla/sessions"
)

type userModule struct {
}

var Module userModule

var DB *sql.DB
var sessionStore sessions.Store

func (module userModule) Name() string {
	return "user"
}

func (module userModule) Version() int {
	return 1
}

func (module userModule) Install() error {
	return nil
}

func (module userModule) Uninstall() error {
	return nil
}

func (module userModule) Upgrade(from, to int) error {
	return nil
}

func (module userModule) SetConnection(conn *sql.DB) {
	DB = conn
}

func (module userModule) SetSessionStore(store sessions.Store) {
	sessionStore = store
}
