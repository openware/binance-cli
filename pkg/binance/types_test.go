package binance

import (
	"encoding/json"
	"testing"

	"github.com/openware/binance-cli/pkg/opendax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinanceMarket(t *testing.T) {
	t.Run("BNBBTC", func(t *testing.T) {
		bm := &BinanceMarket{
			Symbol:         "BNBBTC",
			BaseUnit:       "BNB",
			QuoteUnit:      "BTC",
			QuotePrecision: json.Number("8"),
			Filters: []Filter{
				{
					Type:     "PRICE_FILTER",
					MinPrice: json.Number("0.00000100"),
					MaxPrice: json.Number("100000.00000000"),
					TickSize: json.Number("0.00000100"),
				},
				{
					Type:        "LOT_SIZE",
					MinQuantity: json.Number("0.00100000"),
				},
				{
					Type:        "MIN_NOTIONAL",
					MinNotional: json.Number("0.00010000"),
				},
			},
		}

		odxm, err := bm.ToOpendaxMarket(.0112359550561)
		require.NoError(t, err)

		assert.Equal(t, &opendax.OpendaxMarket{
			Symbol:          "bnbbtc",
			Name:            "BNB/BTC",
			BaseUnit:        "bnb",
			QuoteUnit:       "btc",
			MinPrice:        json.Number("0.00000100"),
			MaxPrice:        json.Number("0.00"),
			MinAmount:       json.Number("0.011"),
			AmountPrecision: 3,
			PricePrecision:  6,
		}, odxm)
	})

	t.Run("BTCUSDT", func(t *testing.T) {
		bm := &BinanceMarket{
			Symbol:         "BTCUSDT",
			BaseUnit:       "BTC",
			QuoteUnit:      "USDT",
			QuotePrecision: json.Number("8"),
			Filters: []Filter{
				{
					Type:     "PRICE_FILTER",
					MinPrice: json.Number("0.01000000"),
					MaxPrice: json.Number("1000000.00000000"),
					TickSize: json.Number("0.01000000"),
				},
				{
					Type:        "LOT_SIZE",
					MinQuantity: json.Number("0.00002000"),
				},
				{
					Type:        "MIN_NOTIONAL",
					MinNotional: json.Number("10.00000000"),
				},
			},
		}

		odxm, err := bm.ToOpendaxMarket(0.0002) // BTC at 50k
		require.NoError(t, err)

		assert.Equal(t, &opendax.OpendaxMarket{
			Symbol:          "btcusdt",
			Name:            "BTC/USDT",
			BaseUnit:        "btc",
			QuoteUnit:       "usdt",
			MinPrice:        json.Number("0.01000000"),
			MaxPrice:        json.Number("0.00"),
			MinAmount:       json.Number("0.00020"),
			AmountPrecision: 5,
			PricePrecision:  2,
		}, odxm)

		odxm, err = bm.ToOpendaxMarket(0.000005) // BTC at 2M
		require.NoError(t, err)

		assert.Equal(t, &opendax.OpendaxMarket{
			Symbol:          "btcusdt",
			Name:            "BTC/USDT",
			BaseUnit:        "btc",
			QuoteUnit:       "usdt",
			MinPrice:        json.Number("0.01000000"),
			MaxPrice:        json.Number("0.00"),
			MinAmount:       json.Number("0.00002000"),
			AmountPrecision: 5,
			PricePrecision:  2,
		}, odxm)

	})
}
