/*
  Package log implements a generic logging system.  Logging can be done on
  different levels (DEBUG, INFO, ...) as well as different aspects of the
  system.  While logging levels are constant, aspects can be created at run
  time, with common aspects being defined globally.
*/
package log

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type logLevel int

// The log levels that are available
const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

var maxLevel = CRITICAL
var levelStrings = [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}

func (l logLevel) String() string {
	if l < 0 || int(l) >= len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

// LogAspect is the type that stores a list of loggers for a given aspect of
// the system.  Each log aspect can be attached to multiple outputs, and a
// logged message will be written to each one.
type LogAspect struct {
	name    string
	loggers map[logLevel][]*log.Logger
}

var aspects map[string]*LogAspect = make(map[string]*LogAspect)

// New creates and returns a named LogAspect.  The created log aspect is stored
// and subsequent calls to New with the same name will return the same object.
// LogAspects are globally defined so that different parts of the system can use
// the same loggers.
func New(name string) *LogAspect {
	key := strings.ToUpper(name)
	if _, ok := aspects[key]; !ok {
		aspect := &LogAspect{
			name:    key,
			loggers: make(map[logLevel][]*log.Logger),
		}
		aspects[key] = aspect
	}
	return aspects[key]
}

var (
	// Database is a default logAspect for logging database events
	Database *LogAspect
	// Datastore is a default LogAspect for datastore events
	Datastore *LogAspect
	all       *LogAspect
)

func init() {
	all = New("All")
	Database = New("Database")
	Datastore = &LogAspect{"ALL", make(map[logLevel][]*log.Logger)}
}

func (aspect LogAspect) prefix(level logLevel) string {
	return fmt.Sprintf("%-8s %-8s ", aspect.name, level)
}

// Add adds an io.Writer for all LogAspects, for levels greater than or equal to logLevel
func Add(baseLevel logLevel, writer io.Writer) {
	all.Add(baseLevel, writer)
}

// AddSingle adds an io.writer for all LogAspects for level
func AddSingle(level logLevel, writer io.Writer) {
	all.AddSingle(level, writer)
}

// AddLogger adds a log.Logger to all LogAspects for levels greater than or
// equal to baseLevel
func AddLogger(baseLevel logLevel, logger *log.Logger) {
	all.AddLogger(baseLevel, logger)
}

// AddSingleLogger adds a log.Logger to all LogAspects, for level
func AddSingleLogger(level logLevel, logger *log.Logger) {
	all.AddSingleLogger(level, logger)
}

// Add adds an io.Writer to a LogAspect, for levels greater than or equal
// to baseLevel
func (aspect LogAspect) Add(baseLevel logLevel, writer io.Writer) {
	for i := int(baseLevel); i <= int(maxLevel); i++ {
		level := logLevel(i)
		logger := log.New(writer, aspect.prefix(level), log.LstdFlags)
		aspect.AddSingleLogger(level, logger)
	}
}

// AddSingle adds an io.Writer for a single log level
func (aspect LogAspect) AddSingle(level logLevel, writer io.Writer) {
	logger := log.New(writer, aspect.prefix(level), log.LstdFlags)
	aspect.loggers[level] = append(aspect.loggers[level], logger)
}

// AddSingleLogger adds a log.Logger for a single log level
func (aspect LogAspect) AddSingleLogger(level logLevel, logger *log.Logger) {
	aspect.loggers[level] = append(aspect.loggers[level], logger)
}

// AddLogger adds a log.Logger for all log levels greater than or equal to
// base level
func (aspect LogAspect) AddLogger(baseLevel logLevel, logger *log.Logger) {
	for i := int(baseLevel); i <= int(maxLevel); i++ {
		aspect.AddSingleLogger(logLevel(i), logger)
	}
}

// Log writes a message to all loggers attached to a log aspect, as well as any
// loggers attached to the global log destination, using a given log level.
func (aspect LogAspect) Log(level logLevel, format string, items ...interface{}) {
	for _, logger := range aspect.loggers[level] {
		logger.Printf(format, items...)
	}
	for _, logger := range all.loggers[level] {
		logger.Printf(format, items...)
	}
}

// Debug writes a message to all attached log destinations, using the DEBUG
// log level
func (aspect LogAspect) Debug(format string, items ...interface{}) {
	aspect.Log(DEBUG, format, items...)
}

// Info writes a message to all attached log destinations, using the INFO
// log level
func (aspect LogAspect) Info(format string, items ...interface{}) {
	aspect.Log(INFO, format, items...)
}

// Warning writes a message to all attached log destinations, using the WARNING
// log level
func (aspect LogAspect) Warning(format string, items ...interface{}) {
	aspect.Log(WARNING, format, items...)
}

// Error writes a message to all attached log destinations, using the ERROR
// log level
func (aspect LogAspect) Error(format string, items ...interface{}) {
	aspect.Log(ERROR, format, items...)
}

// Critical writes a message to all attached log destinations, using the CRITICAL
// log level
func (aspect LogAspect) Critical(format string, items ...interface{}) {
	aspect.Log(CRITICAL, format, items...)
}
