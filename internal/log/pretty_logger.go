package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"sync"
)

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
	w io.Writer
}

const (
	timeFormat = "[15:04:05.000]"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (h *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := fmt.Sprintf("[%s]", r.Level.String())

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	delete(attrs, slog.LevelKey)
	delete(attrs, slog.TimeKey)
	delete(attrs, slog.MessageKey)

	formattedAttrs := ""
	if len(attrs) > 0 {
		bytes, err := json.MarshalIndent(attrs, "", "  ")
		if err != nil {
			return fmt.Errorf("error when marshaling attrs: %w", err)
		}
		formattedAttrs = colorize(darkGray, string(bytes))
	}

	_, err = h.w.Write([]byte(fmt.Sprintln(
		colorize(lightGray, r.Time.Format(timeFormat)),
		level,
		colorize(white, r.Message),
		formattedAttrs,
	)))

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m, w: h.w}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m, w: h.w}
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &Handler{
		b: b,
		h: slog.NewJSONHandler(b, opts),
		m: &sync.Mutex{},
		w: w,
	}
}
