package loglight

import (
	"log"
	"os"
	"strings"
	"runtime"
)

type Logger struct {
	stdout *log.Logger
	packageFilter PackageFilter
}

func NewLogger() *Logger {

	logger := &Logger{
		stdout: log.New(os.Stdout, "",3),
	}

	return logger.WithNoPackageFilter()
}

func (logger *Logger) WithFilter(blacklist PackageFilter) *Logger {
	logger.packageFilter = blacklist
	return logger
}

func (logger *Logger) WithNoPackageFilter() *Logger {
	logger.packageFilter = &NullPackageFilter{}
	return logger
}

func (logger *Logger) LogInfo(msg string) {
	logger.stdout.Print(msg)
}

func (logger *Logger) LogInfof(format string, v ...interface{}) {
	logger.stdout.Printf(format, v)
}

func (logger *Logger) LogDebugf(format string, v ...interface{}) {
	if logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.stdout.Printf(format, v)
	}
}

func (logger *Logger) LogDebug(msg string) {
	if logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.stdout.Printf(msg)
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

