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

type aStruct struct {
	One string
	Two string
}


func TestOutput_NoPackageFilter(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger(true, 0).injectLogPrinter(logPrinter)

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

func TestOutput_NoPackageFilter_NoDebug(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger(false, 0).injectLogPrinter(logPrinter)

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

func TestOutput_WithPackageFilterTrue(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger(true, 0).injectLogPrinter(logPrinter).WithFilter(&mockPackageFilter{filter: true})

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
	logger := NewLogger(true, 0).injectLogPrinter(logPrinter).WithFilter(&mockPackageFilter{filter: false})

	logger.LogDebug("ABC")
	logger.LogInfo("DEF")

	abc, def := false, false

	for _, msg := range logPrinter.messages {
		abc = abc || msg == "ABC"
		def = def || msg == "DEF"
	}

	if abc || def {
		t.Fail()
	}
}


func TestOutput_Struct(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger(true, 0).injectLogPrinter(logPrinter)

	logger.LogInfoStruct(aStruct{One:"1", Two: "2"})

	for _, msg := range logPrinter.messages {
		if msg == "{\"One\":\"1\",\"Two\":\"2\"}" {
			return
		}
	}

	t.Fail()

}


func TestOutput_StructOutputJson(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger(true, 0).injectLogPrinter(logPrinter).OutputJson()

	logger.LogInfoStruct(aStruct{One:"1", Two: "2"})

	for _, msg := range logPrinter.messages {
		if msg == "{\"One\":\"1\",\"Two\":\"2\"}" {
			return
		}
	}

	t.Fail()

}

func Test_OutputJson(t *testing.T) {

	logPrinter := new(mockLogPrinter)
	logger := NewLogger(true, 0).injectLogPrinter(logPrinter).OutputJson()

	logger.LogInfo("test")

	for _, msg := range logPrinter.messages {
		if msg == "{\"Message\":\"test\"}" {
			return
		}
	}

	t.Fail()

}

func TestWhitelist (t *testing.T) {

	filter := NewPackageNameFilter([]string{"test"}, true)

	if !filter.Filter("test") {
		t.Fail()
	}
}

func TestBlacklist (t *testing.T) {

	filter := NewPackageNameFilter([]string{"test"}, false)

	if filter.Filter("test") {
		t.Fail()
	}
}