package log

import (
	"context"
	"log/slog"
	"os"
)

func New(opts ...Option) *slog.Logger {
	config := &Config{}

	for i := range opts {
		opts[i].apply(config)
	}

	if config.writer == nil {
		config.writer = os.Stderr
	}

	var (
		logLevel slog.Level // default is LevelInfo as a zero (int) value
		err      error
	)

	if config.level != "" {
		err = logLevel.UnmarshalText([]byte(config.level))
	}

	if err != nil {
		config.level = "info"
	}

	logger := slog.New(
		NewSpanContextHandler(
			NewPrettyHandler(config.writer, &slog.HandlerOptions{
				AddSource: config.addSource,
				Level:     logLevel,
			}),
			true,
		),
	)

	return logger
}

// NewFromHandler initializes a logger with the provided slog.Handler interface.
//
// If the provided handler is nil, the function returns a no-op slog.Logger.
func NewFromHandler(handler slog.Handler) *slog.Logger {
	if handler == nil {
		return NoOp()
	}

	return slog.New(NewSpanContextHandler(handler, true))
}

func NoOp() *slog.Logger {
	return slog.New(noOpHandler{})
}

type noOpHandler struct {
	slog.Handler
}

func (noOpHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (noOpHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h noOpHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h noOpHandler) WithGroup(string) slog.Handler {
	return h
}
