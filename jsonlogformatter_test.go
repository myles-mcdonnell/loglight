package loglight

import "testing"

func TestFormatStringNotPretty(t *testing.T) {

	var jsonLogFormatter = NewJsonLogFormatter(false)

	var str = jsonLogFormatter.Format(LogEntry{
		LogLevel: DEBUG,
		Data:     "test"})

	var expected = `{"LogLevel":"DEBUG","Data":"test"}`

	if str != expected {

		t.Log("Expected: " + expected + " Actual: " + str)
		t.Fail()
	}
}

func TestFormatStructNotPretty(t *testing.T) {

	var jsonLogFormatter = NewJsonLogFormatter(false)

	var str = jsonLogFormatter.Format(
		LogEntry{
			LogLevel: DEBUG,
			Data: struct {
				Test string
			}{
				Test: "Hello",
			},
		})

	var expected = `{"LogLevel":"DEBUG","Data":{"Test":"Hello"}}`

	if str != expected {

		t.Log("Expected: " + expected + " Actual: " + str)
		t.Fail()
	}
}
