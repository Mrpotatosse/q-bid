package exchange

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllPairs(t *testing.T) {
	ctx := context.Background()

	binanceExchange, err := NewBinanceExchange(
		ctx,
		WithApiKey(os.Getenv(BINANCE_API_KEY)),
		WithSecretKey(os.Getenv(BINANCE_SECRET_KEY)),
		WithAuthorizedCurrencies("BTC", "ETH"),
	)

	require.Nil(t, err)

	pairs := binanceExchange.GetAllPairs("BTC")

	require.Exactly(t, pairs, []*Pair{
		{
			Ask: "ETH",
			Bid: "BTC",
		},
	})
}

func TestGetAllTriples(t *testing.T) {
	ctx := context.Background()

	binanceExchange, err := NewBinanceExchange(
		ctx,
		WithApiKey(os.Getenv(BINANCE_API_KEY)),
		WithSecretKey(os.Getenv(BINANCE_SECRET_KEY)),
		WithAuthorizedCurrencies("BTC", "ETH", "SOL"),
	)

	require.Nil(t, err)

	triples := binanceExchange.GetAllTriples()
	require.Exactly(t, triples, []*Triple{
		{
			End: Pair{
				Ask: "SOL",
				Bid: "BTC",
			},
			Start: Pair{
				Ask: "ETH",
				Bid: "BTC",
			},
		},
		{
			End: Pair{
				Ask: "SOL",
				Bid: "ETH",
			},
			Start: Pair{
				Ask: "ETH",
				Bid: "BTC",
			},
		},
		{
			End: Pair{
				Ask: "SOL",
				Bid: "ETH",
			},
			Start: Pair{
				Ask: "SOL",
				Bid: "BTC",
			},
		},
	})
}
