package loglight

import (
	"log"
	"os"
	"strings"
	"runtime"
)

type packageFilter interface {
	Filter(packageName string) bool
}

type PackageBlacklist struct {
	PackageNames map[string]bool
}

type PackageWhitelist struct {
	PackageNames map[string]bool
}

type NullPackageFilter struct {}

type Logger struct {
	stdout *log.Logger
	packageFilter packageFilter
}

func NewLogger() *Logger {

	logger := &Logger{
		stdout: log.New(os.Stdout, "",3),
	}

	return logger.WithNoPackageFilter()
}

func (logger *Logger) WithBlacklist(blacklist PackageBlacklist) *Logger {
	logger.packageFilter = blacklist
	return logger
}

func (logger *Logger) WithWhitelist(whitelist PackageWhitelist) *Logger {
	logger.packageFilter = whitelist
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

func(blacklist PackageBlacklist) Filter(packageName string) bool {
	return !blacklist.PackageNames[packageName]
}

func(whitelist PackageWhitelist) Filter(packageName string) bool {
	return whitelist.PackageNames[packageName]
}

func(packageFilter NullPackageFilter) Filter(packageName string) bool {
	return true
}