package main

import (
	"github.com/myles-mcdonnell/loglight"
	"github.com/myles-mcdonnell/loglight/example/subPackage"
	"log"
)

func main() {

	defer func() {

		if err := recover(); err != nil {
			log.Fatalf("FATAL ERROR: %s", err)
		}

	} ()

	//all packages other than those listed
	filter := loglight.NewPackageNameFilter([]string{"github.com/myles-mcdonnell/logging/example/subPackage"}, false)
	subPackage.Logger = loglight.NewLogger().WithFilter(filter)

	//No package filter, all debug messages will be written to stdout
	//subPackage.Logger = loglight.NewLogger()

	subPackage.Logger.LogInfo("This message is interesting to users of the software")

	subPackage.Logger.LogDebug("This message provides some information useful to maintainers of the software, such as the internal state")

	subPackage.DoThing()

	panic("Something went horribly wrong that this software can not recover from")

}
