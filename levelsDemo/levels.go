package main

import (
	"fmt"

	"fortio.org/log"
)

func main() {
	log.Config.JSON = true
	log.Config.FatalPanics = false
	log.Config.FatalExit = func(int) {}
	log.SetLogLevelQuiet(log.Debug)
	log.Debugf("This is a debug message ending with backslash \\")
	log.LogVf("This is a verbose message")
	log.Printf("This an always printed, file:line omitted message")
	log.Infof("This is an info message with no attributes but with \"quotes\"...")
	log.S(log.Info, "This is multi line\n\tstructured info message with 3 attributes",
		log.Str("attr1", "value1"), log.Attr("attr2", 42), log.Str("attr3", "\"quoted\nvalue\""))
	log.Warnf("This is a warning message")
	log.Errf("This is an error message")
	log.Critf("This is a critical message")
	log.Fatalf("This is a fatal message") //nolint:revive // we disabled exit for this demo
	fmt.Println("This is a non json output, will get prefixed with a exclamation point")
}
