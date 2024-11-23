package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/Mrpotatosse/q-bid/internal/log"
)

var (
	logLevel = flag.String("log-level", os.Getenv("LOG_LEVEL"), "log level")
)

func main() {
	flag.Parse()
	// Logger configuration
	ctx := context.Background()

	logger := log.New(log.WithLevel(*logLevel))
	ctx = log.WithLogger(ctx, logger)

	if err := run(ctx); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	logger := log.LoggerFromContext(ctx)

	_, err := maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
		logger.Debug(fmt.Sprintf(s, i...))
	}))
	if err != nil {
		return err
	}

	logger.Info("Hello world!", slog.String("location", "world"))

	return nil
}
