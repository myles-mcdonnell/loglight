### loglight

This is a logging package for Go based on Dave Cheney's [blog post about logging](http://dave.cheney.net/2015/11/05/lets-talk-about-logging).

It has all I need for now but could certainly use some more advanced package filtering functionality.  Pull requests are of course very welcome.

Provides ability to customise log output and log entry structures.  Handy when only a subset of field are required for local debugging at development time but full set of log fields are required in production for structured logging, for example.

See [the example command](example/) for usage.

```
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


```



