package subjects

import (
	"database/sql"
	"github.com/YouthBuild-USA/godata/web"
)

type subjectModule struct{}

var db *sql.DB

var Module subjectModule

func (m subjectModule) Name() string {
	return "subjects"
}

func (module subjectModule) Version() int {
	return 1
}

func (module subjectModule) Install() error {
	return nil
}

func (module subjectModule) Uninstall() error {
	return nil
}

func (module subjectModule) Upgrade(from, to int) error {
	return nil
}

func (module subjectModule) SetConnection(conn *sql.DB) {
	db = conn
}

func (module subjectModule) AddRoutes(handleAdder web.HandlerAdder) error {
	subjectTypeRoutes(handleAdder)
	subjectRoutes(handleAdder)
	return nil
}
