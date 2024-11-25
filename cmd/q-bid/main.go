package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	"github.com/Mrpotatosse/q-bid/internal/exchange"
	"github.com/Mrpotatosse/q-bid/internal/log"
	"github.com/joho/godotenv"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	logLevel := flag.String("log-level", os.Getenv("LOG_LEVEL"), "log level")

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

	err = godotenv.Load()
	if err != nil {
		return err
	}

	binanceExchange, err := exchange.NewBinanceExchange(
		ctx,
		exchange.WithApiKey(os.Getenv(exchange.BINANCE_API_KEY)),
		exchange.WithSecretKey(os.Getenv(exchange.BINANCE_SECRET_KEY)),
		exchange.WithAuthorizedCurrencies(
			"BTC",
			"ETH",
			"SOL",
		),
	)
	if err != nil {
		return err
	}

	logger.Debug("test", slog.Any("pairs", binanceExchange.GetAllPairs("BTC")))
	logger.Debug("triples", slog.Any("triples", binanceExchange.GetAllTriples("BTC", "ETH", "SOL")))

	return nil
}
