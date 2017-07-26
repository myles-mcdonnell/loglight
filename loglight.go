package loglight

import (
	"io"
	"os"
	"runtime"
	"strings"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	ERROR          = "ERROR"
	INFO           = "INFO"
)

type LogEntry struct {
	LogLevel
	Data interface{}
}

type GetLogEntry func() interface{}

type LogFormatter func(logEntry LogEntry) string

type Logger struct {
	logWriter     io.Writer
	packageFilter PackageFilter
	outputDebug   bool
	logFormatter  LogFormatter
}

type PackageFilter interface {
	Filter(packageName string) bool
}

type NullPackageFilter struct{}

type packageNameFilter struct {
	packageNames map[string]bool
	isWhitelist  bool
}

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

func NewLogger(outputDebug bool, logFormatter LogFormatter) *Logger {

	return &Logger{
		logWriter:     os.Stdout,
		outputDebug:   outputDebug,
		packageFilter: &NullPackageFilter{},
		logFormatter:  logFormatter,
	}
}

func (logger *Logger) Debug(logEntry interface{}) {
	logger.DebugDefer(func() interface{} { return logEntry })
}

func (logger *Logger) Info(logEntry interface{}) {
	logger.InfoDefer(func() interface{} { return logEntry })
}

func (logger *Logger) Error(logEntry interface{}) {
	logger.ErrorDefer(func() interface{} { return logEntry })
}

func (logger *Logger) DebugDefer(getLogEntryFunc GetLogEntry) {
	logger.writeLogEntry(DEBUG, getLogEntryFunc)
}

func (logger *Logger) InfoDefer(getLogEntryFunc GetLogEntry) {
	logger.writeLogEntry(INFO, getLogEntryFunc)
}

func (logger *Logger) ErrorDefer(getLogEntryFunc GetLogEntry) {
	logger.writeLogEntry(ERROR, getLogEntryFunc)
}

func (logger *Logger) WithLogWriter(logWriter io.Writer) *Logger {
	logger.logWriter = logWriter
	return logger
}

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
