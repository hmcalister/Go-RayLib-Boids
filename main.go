package main

import (
	"flag"
)
func main() {
	rlLogLevelFlag := flag.String("rlLogLevel", "none", "Set the raylib log level. Valid values are: fatal, error, warning, info, debug, trace, none.")
	slogLevelFlag := flag.String("slogLevel", "none", "Set the slog level. Valid values are: fatal, error, warning, info, debug, trace, none.")
	slogFormatFlag := flag.String("slogFormat", "pretty", "Set the slog format. Valid values are: text, pretty, json.")
	flag.Parse()
	setupLogging(*slogLevelFlag, *slogFormatFlag, *rlLogLevelFlag)
}
