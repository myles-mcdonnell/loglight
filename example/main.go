package main

import (
	"fmt"
	"github.com/myles-mcdonnell/loglight"
	"time"
)

func main() {

	logger := loglight.NewLogger(true, loglight.NewJsonLogFormatter(true).Format)

	logger.Info("This message is interesting to users of the software")

	logger.Debug("This message provides some information useful to maintainers of the software, such as the internal state")

	customerLogger := &CustomLogger{logger: loglight.NewLogger(true, CustomFormatter)}

	customerLogger.Debug("123", "custom log event")
}

type CustomLogEvent struct {
	TimeUtc    time.Time
	Key        string
	Additional interface{}
}

func CustomFormatter(logEntry loglight.LogEntry) string {

	logEvent, _ := logEntry.Data.(CustomLogEvent)

	return fmt.Sprintf("%s : %s : %s : %s", logEntry.LogLevel, logEvent.Key, logEvent.TimeUtc.Format("02/01/2006 15:04:05:00"), loglight.GetJson(logEvent.Additional, false))
}

func (logger *CustomLogger) Debug(key string, additional interface{}) {
	logger.logger.Debug(newCustomLogEvent(key, additional))
}

func (logger *CustomLogger) Info(key string, additional interface{}) {
	logger.logger.Info(newCustomLogEvent(key, additional))
}

func (logger *CustomLogger) Error(key string, additional interface{}) {
	logger.logger.Error(newCustomLogEvent(key, additional))
}

func newCustomLogEvent(key string, additional interface{}) CustomLogEvent {
	return CustomLogEvent{
		TimeUtc:    time.Now().UTC(),
		Key:        key,
		Additional: additional,
	}
}

type CustomLogger struct {
	logger *loglight.Logger
}
