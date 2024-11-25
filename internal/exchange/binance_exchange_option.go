package exchange

type BinanceExchangeConfig struct {
	apiKey               string
	secretKey            string
	baseUrl              string
	authorizedCurrencies []string
}

type BinanceExchangeOption interface {
	apply(*BinanceExchangeConfig)
}

type BinanceExchangeOptionFunc func(*BinanceExchangeConfig)

func (o BinanceExchangeOptionFunc) apply(c *BinanceExchangeConfig) {
	o(c)
}

func WithApiKey(key string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.apiKey = key
	}
}

func WithSecretKey(key string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.secretKey = key
	}
}

func WithBaseUrl(url string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.baseUrl = url
	}
}

func WithAuthorizedCurrencies(currencies ...string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.authorizedCurrencies = currencies
	}
}
