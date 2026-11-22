// Package logging is the logger a plugin should log with. gino relays a plugin's stderr into the
// worker's own logs, keeping the level and attributes of JSON lines and logging anything else as
// one info line -- so a plugin logging plain text (slog.Default) loses its levels and attributes.
package logging

import (
	"log/slog"
	"os"

	ginosdkconstants "github.com/gino-works/gino-sdk/constants"
)

// New is a logger writing JSON lines to stderr, at the level gino hands the plugin in
// GinoPluginLogLevelEnv (info when unset or unparseable) -- the worker's own level, so nothing is
// written that gino would only drop.
func New() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: Level(),
			},
		),
	)
}

// Level is the level gino hands the plugin in GinoPluginLogLevelEnv, info when unset or
// unparseable.
func Level() slog.Level {
	var level slog.Level

	err := level.UnmarshalText([]byte(os.Getenv(ginosdkconstants.GinoPluginLogLevelEnv)))
	if err != nil {
		return slog.LevelInfo
	}

	return level
}
