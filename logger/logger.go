package logger

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/lmittmann/tint"
)

func SetGlobalLogger() {
	minLevel, err := strconv.Atoi(os.Getenv("LOGGER_LEVEL"))
	if err != nil {
		minLevel = int(slog.LevelInfo)
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.Level(minLevel),
			TimeFormat: time.Kitchen,
		}),
	))
}
