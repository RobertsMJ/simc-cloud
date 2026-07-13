package applog

import (
	"log/slog"
	"os"
)

func Init() {
	logLevel := slog.LevelInfo

	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		_ = logLevel.UnmarshalText([]byte(envLevel))
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
}
