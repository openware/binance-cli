package binance

func NewBinanceClient(apiKey, secret string) *BinanceClient {
	return &BinanceClient{
		apiKey: apiKey,
		secret: secret,
	}
}

func (bc *BinanceClient) CoinsInfo() (BinanceCurrencies, error) {
	currencies := BinanceCurrencies{}
	_, err := bc.apiCall(coinsInfoEndpoint, &currencies)
	return currencies, err
}

func (bc *BinanceClient) ExchangeInfo() (BinanceExchangeInfo, error) {
	exchangeInfo := BinanceExchangeInfo{}
	_, err := bc.apiCall(exchangeInfoEndpoint, &exchangeInfo)
	exchangeInfo.FillRegistry()
	return exchangeInfo, err
}
