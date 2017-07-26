package loglight

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
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

	buf := new(bytes.Buffer)
	logger := NewLogger(true, NewJsonLogFormatter(false).Format).WithLogWriter(buf)

	logger.DebugDefer(func() interface{} { return "ABC" })
	logger.InfoDefer(func() interface{} { return "DEF" })

	reader := bufio.NewReader(buf)

	CompareString(reader, `{"LogLevel":"DEBUG","Data":"ABC"}`, t)
	CompareString(reader, `{"LogLevel":"INFO","Data":"DEF"}`, t)
}

func CompareString(reader *bufio.Reader, expected string, t *testing.T) {
	str, err := reader.ReadString('\n')
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if strings.Compare(strings.Trim(str, "\n"), expected) != 0 {
		t.Log("Expected: " + expected + " Actual: " + str)
		t.Fail()
	}
}

func TestOutput_NoPackageFilter_NoDebug(t *testing.T) {

	buf := new(bytes.Buffer)
	logger := NewLogger(false, NewJsonLogFormatter(false).Format).WithLogWriter(buf)

	logger.DebugDefer(func() interface{} { return "ABC" })
	logger.InfoDefer(func() interface{} { return "DEF" })

	reader := bufio.NewReader(buf)

	CompareString(reader, `{"LogLevel":"INFO","Data":"DEF"}`, t)
}

func TestOutput_WithPackageFilterTrue(t *testing.T) {

	buf := new(bytes.Buffer)
	logger := NewLogger(true, NewJsonLogFormatter(false).Format).WithLogWriter(buf).WithFilter(&mockPackageFilter{filter: true})

	logger.DebugDefer(func() interface{} { return "ABC" })
	logger.InfoDefer(func() interface{} { return "DEF" })

	reader := bufio.NewReader(buf)

	CompareString(reader, `{"LogLevel":"DEBUG","Data":"ABC"}`, t)
	CompareString(reader, `{"LogLevel":"INFO","Data":"DEF"}`, t)
}

func TestOutput_WithPackageFilterFalse(t *testing.T) {

	buf := new(bytes.Buffer)
	logger := NewLogger(true, NewJsonLogFormatter(false).Format).WithLogWriter(buf).WithFilter(&mockPackageFilter{filter: false})

	//logger.Debug(func() interface{} { return "ABC" })
	//logger.Info(func() interface{} { return "DEF" })

	logger.Debug("ABC")
	logger.Info("DEF")

	reader := bufio.NewReader(buf)

	_, err := reader.ReadString('\n')

	if err != io.EOF {
		t.Fail()
	}
}

func TestWhitelist(t *testing.T) {

	filter := NewPackageNameFilter([]string{"test"}, true)

	if !filter.Filter("test") {
		t.Fail()
	}
}

func TestBlacklist(t *testing.T) {

	filter := NewPackageNameFilter([]string{"test"}, false)

	if filter.Filter("test") {
		t.Fail()
	}
}
