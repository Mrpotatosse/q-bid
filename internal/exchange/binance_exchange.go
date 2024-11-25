package exchange

import (
	"context"

	binance "github.com/binance/binance-connector-go"
)

type BinanceExchange struct {
	Exchange

	ctx    context.Context
	client *binance.Client

	pairs []*Pair
}

func NewBinanceExchange(ctx context.Context, opts ...BinanceExchangeOption) (*BinanceExchange, error) {
	config := &BinanceExchangeConfig{}

	for _, opt := range opts {
		opt.apply(config)
	}

	var client *binance.Client
	if config.baseUrl != "" {
		client = binance.NewClient(config.apiKey, config.secretKey, config.baseUrl)
	} else {
		client = binance.NewClient(config.apiKey, config.secretKey)
	}

	exchangeInfoService, err := client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return nil, err
	}

	pairs := []*Pair{}
	pairsMap := make(map[string][]string)
	for _, symbol := range exchangeInfoService.Symbols {
		pairsMap[symbol.BaseAsset] = append(pairsMap[symbol.BaseAsset], symbol.QuoteAsset)
	}

	for _, currencyAsk := range config.authorizedCurrencies {
		pairsList, found := pairsMap[currencyAsk]

		if found {
			for _, currencyBid := range pairsList {
				for _, authorizedCurrency := range config.authorizedCurrencies {
					if currencyBid == authorizedCurrency {
						pairs = append(pairs, &Pair{
							Ask: currencyAsk,
							Bid: currencyBid,
						})
						break
					}
				}

			}
		}
	}

	return &BinanceExchange{
		ctx:    ctx,
		client: client,

		pairs: pairs,
	}, nil
}

func (e *BinanceExchange) GetAllPairs(currencies ...string) (pairs []*Pair) {
	if len(currencies) == 0 {
		return e.pairs
	}

	pairs = []*Pair{}

	for _, pair := range e.pairs {
		for _, currency := range currencies {
			if pair.Ask == currency || pair.Bid == currency {
				pairs = append(pairs, pair)
				break
			}
		}
	}

	return
}

func (e *BinanceExchange) GetAllTriples(currencies ...string) (triples []*Triple) {
	triples = []*Triple{}
	pairs := e.GetAllPairs(currencies...)

	for i, startPair := range pairs {
		for _, endPair := range pairs[i+1:] {
			if startPair.Ask == endPair.Ask ||
				startPair.Ask == endPair.Bid ||
				startPair.Bid == endPair.Ask ||
				startPair.Bid == endPair.Bid {
				triples = append(triples, &Triple{
					Start: *startPair,
					End:   *endPair,
				})
			}
		}
	}

	return
}
