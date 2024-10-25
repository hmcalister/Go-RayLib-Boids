package main

import (
	"flag"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rlLogLevelFlag := flag.String("rlLogLevel", "none", "Set the raylib log level. Valid values are: fatal, error, warning, info, debug, trace, none.")
	slogLevelFlag := flag.String("slogLevel", "info", "Set the slog level. Valid values are: fatal, error, warning, info, debug, none.")
	slogFormatFlag := flag.String("slogFormat", "pretty", "Set the slog format. Valid values are: text, pretty, json.")
	configFilePath := flag.String("configFilePath", "config.yaml", "Path to the config file to use.")
	flag.Parse()
	setupLogging(*slogLevelFlag, *slogFormatFlag, *rlLogLevelFlag)

	config, err := ParseConfigFile(*configFilePath)
	if err != nil {
		os.Exit(1)
	}

	rl.InitWindow(config.WindowWidth, config.WindowHeight, "Boids")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	manager := NewBoidManager(config)

	for !rl.WindowShouldClose() {
		rl.ClearBackground(rl.Black)

		for _, b := range manager.Boids {
			b.DrawBoid()
		}

		rl.EndDrawing()
	}

}
