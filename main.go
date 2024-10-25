package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	rlLogLevelFlag := flag.String("rlLogLevel", "none", "Set the raylib log level. Valid values are: fatal, error, warning, info, debug, trace, none.")
	slogLevelFlag := flag.String("slogLevel", "none", "Set the slog level. Valid values are: fatal, error, warning, info, debug, trace, none.")
	slogFormatFlag := flag.String("slogFormat", "pretty", "Set the slog format. Valid values are: text, pretty, json.")
	configFilePath := flag.String("configFilePath", "config.yaml", "Path to the config file to use.")
	flag.Parse()
	setupLogging(*slogLevelFlag, *slogFormatFlag, *rlLogLevelFlag)

	config, err := ParseConfigFile(*configFilePath)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%+v\n", config)
}
