package loglight

import (
	"io"
	"os"
	"runtime"
	"strings"
)

type LogLevel string

// Loglight has only three log levels
const (
	DEBUG LogLevel = "DEBUG"
	ERROR          = "ERROR"
	INFO           = "INFO"
)

// The default LogEntry type
type LogEntry struct {
	LogLevel
	Data interface{}
}

// Construction of LogEntries may be deferred until the log entry is written.
// This is a performance enhancement to prevent construction of Debug log entries that may not be required when Debug output is disabled
type GetLogEntry func() interface{}

// The function used to format the log output.  Custom implementations may be supplied to the Logger constructor NewLogger
type LogFormatter func(logEntry LogEntry) string

// The Logger on which all log methods are defined
type Logger struct {
	logWriter     io.Writer
	packageFilter PackageFilter
	outputDebug   bool
	logFormatter  LogFormatter
}

// Used to selectively include or exclude packages from the log output
type PackageFilter interface {
	Filter(packageName string) bool
}

type NullPackageFilter struct{}

type packageNameFilter struct {
	packageNames map[string]bool
	isWhitelist  bool
}

// Creates a package filter that includes or excludes a list of packages
func NewPackageNameFilter(packageNames []string, isWhitelist bool) PackageFilter {

	filter := &packageNameFilter{isWhitelist: isWhitelist, packageNames: make(map[string]bool)}

	for _, name := range packageNames {
		filter.packageNames[name] = true
	}

	return filter
}

func (filter packageNameFilter) Filter(packageName string) bool {
	return filter.packageNames[packageName] && filter.isWhitelist
}

func (packageFilter NullPackageFilter) Filter(packageName string) bool {
	return true
}

// Logger constructor
func NewLogger(outputDebug bool, logFormatter LogFormatter) *Logger {

	return &Logger{
		logWriter:     os.Stdout,
		outputDebug:   outputDebug,
		packageFilter: &NullPackageFilter{},
		logFormatter:  logFormatter,
	}
}

// Log Debug event
func (logger *Logger) Debug(logEntry interface{}) {
	logger.DebugDefer(func() interface{} { return logEntry })
}

// Log Info event
func (logger *Logger) Info(logEntry interface{}) {
	logger.InfoDefer(func() interface{} { return logEntry })
}

// Log Error event
func (logger *Logger) Error(logEntry interface{}) {
	logger.ErrorDefer(func() interface{} { return logEntry })
}

// Log Debug event and defer log entry construction
func (logger *Logger) DebugDefer(getLogEntryFunc GetLogEntry) {
	logger.writeLogEntry(DEBUG, getLogEntryFunc)
}

// Log Info event and defer log entry construction
func (logger *Logger) InfoDefer(getLogEntryFunc GetLogEntry) {
	logger.writeLogEntry(INFO, getLogEntryFunc)
}

// Log Error event and defer log entry construction
func (logger *Logger) ErrorDefer(getLogEntryFunc GetLogEntry) {
	logger.writeLogEntry(ERROR, getLogEntryFunc)
}

// Inject an alternative io.Writer for log output
func (logger *Logger) WithLogWriter(logWriter io.Writer) *Logger {
	logger.logWriter = logWriter
	return logger
}

// Inject log filter
func (logger *Logger) WithFilter(filter PackageFilter) *Logger {
	logger.packageFilter = filter
	return logger
}

func (logger *Logger) writeLogEntry(level LogLevel, getLogEntryFunc GetLogEntry) {

	if (level != DEBUG || logger.outputDebug) && logger.packageFilter.Filter(retrieveCallPackage()) {

		logEntry := LogEntry{
			LogLevel: level,
			Data:     getLogEntryFunc(),
		}

		io.WriteString(logger.logWriter, logger.logFormatter(logEntry)+"\n")
	}
}

func retrieveCallPackage() string {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""

	if parts[pl-2][0] == '(' {
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return packageName
}
