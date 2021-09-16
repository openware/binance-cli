package binance

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

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
			MinPrice: decimal.RequireFromString("0.01000000"),
			MaxPrice: decimal.RequireFromString("1000000.00000000"),
			TickSize: decimal.RequireFromString("0.01000000"),
		},
		{
			Type:        "LOT_SIZE",
			MinQuantity: decimal.RequireFromString("0.00010000"),
		},
		{
			Type:        "MIN_NOTIONAL",
			MinNotional: decimal.RequireFromString("10.00000000"),
		},
	}

	expectedMarket = BinanceMarket{
		Symbol:         "ETHUSDT",
		BaseUnit:       "ETH",
		QuoteUnit:      "USDT",
		QuotePrecision: decimal.RequireFromString("8"),
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
	assert.NilError(t, err)

	expectedTickerPriceRes := &BinanceTickerPrice{
		Symbol: "ETHUSDT",
		Price:  decimal.RequireFromString("3500.00000000"),
	}

	assert.DeepEqual(t, expectedTickerPriceRes, res)
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
	assert.NilError(t, err)

	expectedExchangeInfoRes := &BinanceExchangeInfo{
		Symbols:        []BinanceMarket{expectedMarket},
		MarketRegistry: map[string]BinanceMarket{"ETHUSDT": expectedMarket},
	}

	assert.DeepEqual(t, expectedExchangeInfoRes, res)
}

func TestCalculateMinAmount(t *testing.T) {
	expectedTickerPriceRes := &BinanceTickerPrice{
		Symbol: "ETHUSDT",
		Price:  decimal.RequireFromString("3500.00000000"),
	}

	minAmount := expectedMarket.CalculateMinAmount(expectedTickerPriceRes.Price)

	assert.DeepEqual(t, decimal.RequireFromString("0.003"), minAmount)
}

func TestToOpendaxMarket(t *testing.T) {
	expectedOpendaxMarket := &opendax.OpendaxMarket{
		Symbol:          "ethusdt",
		Name:            "ETH/USDT",
		BaseUnit:        "eth",
		QuoteUnit:       "usdt",
		MinPrice:        expectedFilters[0].MinPrice,
		MaxPrice:        decimal.RequireFromString("0.00"),
		MinAmount:       decimal.RequireFromString("0.0030"),
		AmountPrecision: int64(4),
		PricePrecision:  int64(2),
	}

	res, err := expectedMarket.ToOpendaxMarket(decimal.RequireFromString("0.003"))
	require.NoError(t, err)
	assert.DeepEqual(t, expectedOpendaxMarket, res)
}
