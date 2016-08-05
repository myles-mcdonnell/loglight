package loglight

import (
	"fmt"
	"testing"
)

type mockLogPrinter struct {
	messages []string
}

type mockPackageFilter struct {
	filter bool
}

func (mockLogPrinter *mockLogPrinter) Printf(format string, v ...interface{}) {
	mockLogPrinter.messages = append(mockLogPrinter.messages, fmt.Sprintf(format, v...))
}

func (mockLogPrinter *mockLogPrinter) Print(v ...interface{}) {
	mockLogPrinter.messages = append(mockLogPrinter.messages, fmt.Sprint(v...))
}

func (mockPackageFilter *mockPackageFilter) Filter(packageName string) bool {
	return mockPackageFilter.filter
}


func TestOutput_NoPackageFilter(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger().injectLogPrinter(logPrinter)

	logger.LogDebug("ABC")
	logger.LogInfo("DEF")

	abc, def := false, false

	for _, msg := range logPrinter.messages {
		abc = abc || msg == "ABC"
		def = def || msg == "DEF"
	}

	if !abc || !def {
		t.Fail()
	}
}


func TestOutput_WithPackageFilterTrue(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger().injectLogPrinter(logPrinter).WithFilter(&mockPackageFilter{filter: true})

	logger.LogDebug("ABC")
	logger.LogInfo("DEF")

	abc, def := false, false

	for _, msg := range logPrinter.messages {
		abc = abc || msg == "ABC"
		def = def || msg == "DEF"
	}

	if !abc || !def {
		t.Fail()
	}
}

func TestOutput_WithPackageFilterFalse(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger().injectLogPrinter(logPrinter).WithFilter(&mockPackageFilter{filter: false})

	logger.LogDebug("ABC")
	logger.LogInfo("DEF")

	abc, def := false, false

	for _, msg := range logPrinter.messages {
		abc = abc || msg == "ABC"
		def = def || msg == "DEF"
	}

	if abc || !def {
		t.Fail()
	}
}