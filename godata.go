package godata

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/YouthBuild-USA/godata/config"
	"github.com/YouthBuild-USA/godata/global"
	"github.com/YouthBuild-USA/godata/log"
	"github.com/YouthBuild-USA/godata/questions"
	"github.com/YouthBuild-USA/godata/subjects"
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/users"
	"github.com/YouthBuild-USA/godata/web"
	_ "github.com/bmizerany/pq"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	golog "log"
	"net/http"
)

var db *sql.DB = nil
var dbMap *gorp.DbMap

var AssetDirectory string = "assets"
var SessionKey string = "this should be changed!"

var ModuleLog *log.LogAspect
var DBLog *log.LogAspect

func init() {
	ModuleLog = log.New("Module")
	DBLog = log.New("Database")
	Register(users.Module)

	config.Register("DEFAULT", "port", "8080", "The port on which to run the webserver")
	config.Register("DEFAULT", "assetDirectory", "assets", `
		The directory to find static assets and templates. Can be an absolute path or
		relative to the executable.
		`)

	config.Register("Database", "username", "", "The database username")
	config.Register("Database", "password", "", "The database password")
	config.Register("Database", "database", "", "The name of the database to use")
	config.Register("Database", "host", "localhost", "The database host")
}

var modules map[string]Module = make(map[string]Module)

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
func Start(configFile string) (chan error, error) {
	config.SetFile(configFile)

	dbUser := config.MustGet("Database", "username")
	dbPass := config.MustGet("Database", "password")
	dbName := config.MustGet("Database", "database")
	dbHost := config.MustGet("Database", "host")

	db, _ = sql.Open("postgres", fmt.Sprintf("user=%v password=%v dbname=%v host=%v", dbUser, dbPass, dbName, dbHost))
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	global.DB = db
	global.DbMap = dbMap

	dbLogger := golog.New(DBLog.Writer(log.INFO), "", 0)
	dbMap.TraceOn("GORP", dbLogger)

	criticalErrors := make(chan error)

	web.CreateSessionStore(SessionKey)

	if db == nil {
		return nil, errors.New("Database connection must be set before call to Start")
	}

	router := mux.NewRouter()

	templates.Router = router

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

	questions.InitializeDatabase(db, dbMap)
	subjects.Initialize(dbMap)

	questions.AddRoutes(adder)
	subjects.AddRoutes(adder)

	http.Handle("/static/", http.FileServer(http.Dir(AssetDirectory)))
	http.Handle("/", router)
	port, _ := config.Get("DEFAULT", "port")
	http.ListenAndServe(":"+port, nil)

	return criticalErrors, nil
}

// bootstrapModule runs any required steps to initialize modules
func boostrapModule(module Module) error {
	if dbModule, ok := module.(DatabaseModule); ok {
		dbModule.SetConnection(db)
	}
	return nil
}
