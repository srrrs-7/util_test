package logger

import (
	"log/slog"
	"os"
)

func New() {
	// slog json logger
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
