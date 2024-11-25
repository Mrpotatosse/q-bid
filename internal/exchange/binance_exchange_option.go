package exchange

type BinanceExchangeConfig struct {
	apiKey               string
	secretKey            string
	baseURL              string
	authorizedCurrencies []string
}

type BinanceExchangeOption interface {
	apply(*BinanceExchangeConfig)
}

type BinanceExchangeOptionFunc func(*BinanceExchangeConfig)

func (o BinanceExchangeOptionFunc) apply(c *BinanceExchangeConfig) {
	o(c)
}

func WithAPIKey(key string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.apiKey = key
	}
}

func WithSecretKey(key string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.secretKey = key
	}
}

func WithBaseURL(url string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.baseURL = url
	}
}

func WithAuthorizedCurrencies(currencies ...string) BinanceExchangeOptionFunc {
	return func(bec *BinanceExchangeConfig) {
		bec.authorizedCurrencies = currencies
	}
}
