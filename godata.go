package godata

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/YouthBuild-USA/godata/log"
	"github.com/YouthBuild-USA/godata/users"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/mux"
	"net/http"
)

var db *sql.DB = nil

var AssetDirectory string = "assets"
var SessionKey string = "this should be changed!"

var ModuleLog *log.LogAspect

func init() {
	ModuleLog = log.New("Module")

	Register(users.Module)
}

var modules map[string]Module = make(map[string]Module)

// SetDatabase sets the database connection to use for the system
func SetDatabase(connection *sql.DB) {
	db = connection
}

// Register registers a module with the system
func Register(module Module) error {
	if _, ok := modules[module.Name()]; ok {
		return ModuleExistsError(module.Name())
	}

	modules[module.Name()] = module
	return nil
}

// Start starts the system.  If the second error parameter is not nil, then the
// system failed to start.  Future system halting errors are captured and
// returned on the channel.
func Start() (chan error, error) {
	criticalErrors := make(chan error)

	web.CreateSessionStore(SessionKey)

	if db == nil {
		return nil, errors.New("Database connection must be set before call to Start")
	}

	router := mux.NewRouter()

	adder := func(path string, handle web.Handle) *mux.Route {
		ModuleLog.Info("Registered Path %v", path)
		return router.Handle(path, handle)
	}

	for _, module := range modules {
		err := boostrapModule(module)
		if err != nil {
			return nil, err
		}

		if pageModule, ok := module.(PageModule); ok {
			ModuleLog.Info(fmt.Sprintf("Registering paths from %v", module.Name()))
			pageModule.AddRoutes(adder)
		}
	}

	http.Handle("/static/", http.FileServer(http.Dir(AssetDirectory)))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

	return criticalErrors, nil
}

// bootstrapModule runs any required steps to initialize modules
func boostrapModule(module Module) error {
	if dbModule, ok := module.(DatabaseModule); ok {
		dbModule.SetConnection(db)
	}
	return nil
}
