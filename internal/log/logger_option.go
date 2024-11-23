package log

import (
	"context"
	"io"
)

type Config struct {
	writer    io.Writer
	addSource bool
	level     string
	context   context.Context
}

type Option interface {
	apply(*Config)
}

type OptionFunc func(*Config)

func (o OptionFunc) apply(c *Config) {
	o(c)
}

func WithWriter(writer io.Writer) OptionFunc {
	if writer == nil {
		return func(*Config) {}
	}

	return func(config *Config) {
		config.writer = writer
	}
}

func WithLevel(level string) OptionFunc {
	if level == "" {
		return func(*Config) {}
	}

	return func(config *Config) {
		config.level = level
	}
}

func WithSource() OptionFunc {
	return func(config *Config) {
		config.addSource = true
	}
}

func WithContext(ctx context.Context) OptionFunc {
	return func(config *Config) {
		config.context = ctx
	}
}
