package loglight

import (
	"fmt"
	"gopkg.in/myles-mcdonnell/jsonx.v1"
)

type message struct {
	Message string
}

type JsonLogFormatter struct {
	pretty bool
}

func NewJsonLogFormatter(pretty bool) JsonLogFormatter {
	return JsonLogFormatter{pretty: pretty}
}

func (jsonLogFormatter JsonLogFormatter) Format(logEntry LogEntry) string {

	return getJson(logEntry, jsonLogFormatter.pretty)
}

func getJson(msg interface{}, pretty bool) string {

	bytes, err := marshallJson(msg, pretty)

	if err != nil {
		return fmt.Sprintf("error serializing msg %s", err.Error())
	}

	return string(bytes)
}

func marshallJson(msg interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return jsonx.MarshalIndentWithOptions(msg, "", "    ", jsonx.MarshalOptions{SkipUnserializableFields: true})
	} else {
		return jsonx.MarshalWithOptions(msg, jsonx.MarshalOptions{SkipUnserializableFields: true})
	}
}
