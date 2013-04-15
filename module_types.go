package godata

import (
	"database/sql"
	"github.com/YouthBuild-USA/godata/web"
)

type Module interface {
	Name() string
}

type DatabaseModule interface {
	SetConnection(db *sql.DB)
}

type PageModule interface {
	AddRoutes(handlerAdder web.HandlerAdder) error
}
