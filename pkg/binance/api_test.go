package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openware/binance-cli/pkg/opendax"
)

var (
	mux           *http.ServeMux
	server        *httptest.Server
	binanceClient *BinanceClient
	ticker        = "ETHUSDT"

	expectedFilters = []Filter{
		{
			Type:     "PRICE_FILTER",
			MinPrice: json.Number("0.01000000"),
			MaxPrice: json.Number("1000000.00000000"),
			TickSize: json.Number("0.01000000"),
		},
		{
			Type:        "LOT_SIZE",
			MinQuantity: json.Number("0.00010000"),
		},
		{
			Type:        "MIN_NOTIONAL",
			MinNotional: json.Number("10.00000000"),
		},
	}

	expectedMarket = BinanceMarket{
		Symbol:         "ETHUSDT",
		BaseUnit:       "ETH",
		QuoteUnit:      "USDT",
		QuotePrecision: json.Number("8"),
		Filters:        expectedFilters,
	}
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	binanceClient = &BinanceClient{
		apiKey: "tr13dt0",
		secret: "f0rg3t",
		url:    server.URL,
	}

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/" + path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func TestTickerPriceEndpoint(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(tickerPriceInfoEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("ticker_price_ethusdt.json"))
	})

	res, err := binanceClient.TickerPriceInfo(ticker)
	assert.NoError(t, err)

	expectedTickerPriceRes := &BinanceTickerPrice{
		Symbol: "ETHUSDT",
		Price:  json.Number("3500.00000000"),
	}

	assert.Equal(t, expectedTickerPriceRes, res)
}

func TestExchangeInfoEndpoint(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(exchangeInfoEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("exchange_info.json"))
	})

	res, err := binanceClient.ExchangeInfo()
	assert.NoError(t, err)

	expectedExchangeInfoRes := &BinanceExchangeInfo{
		Symbols:        []BinanceMarket{expectedMarket},
		MarketRegistry: map[string]BinanceMarket{"ETHUSDT": expectedMarket},
	}

	assert.Equal(t, expectedExchangeInfoRes, res)
}

func TestCalculateMinAmount(t *testing.T) {
	expectedTickerPriceRes := &BinanceTickerPrice{
		Symbol: "ETHUSDT",
		Price:  json.Number("3500.00000000"),
	}

	minAmount, err := expectedMarket.CalculateMinAmount(expectedTickerPriceRes.Price)
	assert.NoError(t, err)

	assert.Equal(t, 0.003, minAmount)
}

func TestToOpendaxMarket(t *testing.T) {
	expectedOpendaxMarket := &opendax.OpendaxMarket{
		Symbol:          "ethusdt",
		Name:            "ETH/USDT",
		BaseUnit:        "eth",
		QuoteUnit:       "usdt",
		MinPrice:        expectedFilters[0].MinPrice,
		MaxPrice:        json.Number("0.00"),
		MinAmount:       json.Number("0.0030"),
		AmountPrecision: 4,
		PricePrecision:  2,
	}

	res, err := expectedMarket.ToOpendaxMarket(0.003)
	require.NoError(t, err)
	assert.Equal(t, expectedOpendaxMarket, res)
}
