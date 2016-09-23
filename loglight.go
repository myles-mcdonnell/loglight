package loglight

import (
	"log"
	"os"
	"strings"
	"fmt"
	"runtime"
	"gopkg.in/myles-mcdonnell/jsonx.v1"
)

type logPrinter interface {
	Print(...interface{})
	Printf(string, ...interface{})
}

type PackageFilter interface {
	Filter(packageName string) bool
}

type Logger struct {
	logPrinter logPrinter
	packageFilter PackageFilter
	outputDebug bool
}

type NullPackageFilter struct {}

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

func(filter packageNameFilter) Filter(packageName string) bool {
	return filter.packageNames[packageName] && filter.isWhitelist
}

func(packageFilter NullPackageFilter) Filter(packageName string) bool {
	return true
}

func NewLogger(outputDebug bool) *Logger {

	return NewLoggerWithFlag(outputDebug, 3)
}

func NewLoggerWithFlag(outputDebug bool, flag int) *Logger {

	logger := &Logger{
		logPrinter: log.New(os.Stdout, "",flag),
		outputDebug: outputDebug,
	}

	return logger.WithNoPackageFilter()
}

func (logger *Logger) WithFilter(filter PackageFilter) *Logger {
	logger.packageFilter = filter
	return logger
}

func (logger *Logger) WithNoPackageFilter() *Logger {
	logger.packageFilter = &NullPackageFilter{}
	return logger
}

func (logger *Logger) injectLogPrinter(logPrinter logPrinter) *Logger {

	logger.logPrinter = logPrinter

	return logger
}

func (logger *Logger) LogInfo(msg string) {
	if logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.logPrinter.Print(msg)
	}
}

func (logger *Logger) LogInfof(format string, v ...interface{}) {
	if logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.logPrinter.Printf(format, v)
	}
}

func (logger *Logger) LogInfoStruct(msg interface{}) {
	if logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.logPrinter.Print(getJson(msg))
	}
}

func (logger *Logger) LogDebugf(format string, v ...interface{}) {
	if logger.outputDebug && logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.logPrinter.Printf(format, v)
	}
}

func (logger *Logger) LogDebug(msg string) {
	if logger.outputDebug && logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.logPrinter.Print(msg)
	}
}

func (logger *Logger) LogDebugStruct(msg interface{}) {
	if logger.outputDebug && logger.packageFilter.Filter(retrieveCallPackage()) {
		logger.logPrinter.Print(getJson(msg))
	}
}

func getJson(msg interface{}) string {

	bytes, err := jsonx.MarshalWithOptions(msg, jsonx.MarshalOptions{SkipUnserializableFields:true})

	if err != nil {
		return fmt.Sprintf("error serializing msg %s", err.Error())
	}

	return string(bytes)
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

