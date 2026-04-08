package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New returns a JSON slog logger configured for the given level and optional environment label.
func New(level, env string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	l := slog.New(handler).With("env", env)
	slog.SetDefault(l)
	return l
}
