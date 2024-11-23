package log

import (
	"context"
	"log/slog"
)

// This code is used to access logger from a context
type LoggerKey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey{}, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerKey{}).(*slog.Logger)

	if !ok {
		return slog.Default()
	}

	return logger
}
