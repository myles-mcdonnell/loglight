package subPackage

import (
	"github.com/myles-mcdonnell/loglight"
)

var Logger *loglight.Logger

func DoThing() {

	Logger.LogDebug("You may or may not be interested in this message")
}