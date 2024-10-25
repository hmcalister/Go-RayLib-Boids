package main

import (
	"log/slog"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	console "github.com/phsym/console-slog"
)

func setupLogging(slogLevelString string, slogFormatString string, rlLevelString string) {
	var slogLevel slog.Level
	switch slogLevelString {
	case "error":
		slogLevel = slog.LevelError
	case "warn":
		slogLevel = slog.LevelWarn
	case "info":
		slogLevel = slog.LevelInfo
	case "debug":
		slogLevel = slog.LevelDebug
	case "none":
		slogLevel = slog.LevelError + 1
	default:
		slog.Error("flag is not a valid slog level", "slog level flag", slogLevelString)
		os.Exit(1)
	}

	var slogHandler slog.Handler
	switch slogFormatString {
	case "text":
		slogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		})
	case "pretty":
		slogHandler = console.NewHandler(os.Stdout, &console.HandlerOptions{
			Level: slogLevel,
		})
	case "json":
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slogLevel,
		})
	default:
		slog.Error("flag is not a valid slog format", "slog format flag", slogFormatString)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(
		slogHandler,
	))

	// --------------------------------------------------------------------------------

	var rlLogLevel rl.TraceLogLevel
	switch rlLevelString {
	case "fatal":
		rlLogLevel = rl.LogFatal
	case "error":
		rlLogLevel = rl.LogError
	case "warning":
		rlLogLevel = rl.LogWarning
	case "info":
		rlLogLevel = rl.LogInfo
	case "debug":
		rlLogLevel = rl.LogDebug
	case "trace":
		rlLogLevel = rl.LogTrace
	case "none":
		rlLogLevel = rl.LogNone
	default:
		slog.Error("flag is not a valid raylib log level", "raylib log level flag", rlLevelString)
		os.Exit(1)
	}

	rl.SetTraceLogLevel(rlLogLevel)
}
