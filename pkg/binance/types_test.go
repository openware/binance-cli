package binance

import (
	"fmt"
	"testing"

	"github.com/openware/binance-cli/pkg/opendax"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestJsonNumbersEqual(t *testing.T) {
	assert.Equal(t, true, decimal.RequireFromString("11").Equals(decimal.RequireFromString("11.0")))
}

func TestBinanceMarket(t *testing.T) {
	t.Run("BNBBTC", func(t *testing.T) {
		bm := &BinanceMarket{
			Symbol:         "BNBBTC",
			BaseUnit:       "BNB",
			QuoteUnit:      "BTC",
			QuotePrecision: decimal.RequireFromString("8"),
			Filters: []Filter{
				{
					Type:     "PRICE_FILTER",
					MinPrice: decimal.RequireFromString("0.00000100"),
					MaxPrice: decimal.RequireFromString("100000.00000000"),
					TickSize: decimal.RequireFromString("0.00000100"),
				},
				{
					Type:        "LOT_SIZE",
					MinQuantity: decimal.RequireFromString("0.00100000"),
				},
				{
					Type:        "MIN_NOTIONAL",
					MinNotional: decimal.RequireFromString("0.00010000"),
				},
			},
		}

		odxm, err := bm.ToOpendaxMarket(decimal.NewFromFloat32(float32(0.0112359550561)))
		require.NoError(t, err)

		assert.DeepEqual(t, &opendax.OpendaxMarket{
			Symbol:          "bnbbtc",
			Name:            "BNB/BTC",
			BaseUnit:        "bnb",
			QuoteUnit:       "btc",
			MinPrice:        decimal.RequireFromString("0.00000100"),
			MaxPrice:        decimal.RequireFromString("0.00"),
			MinAmount:       decimal.RequireFromString("0.011"),
			AmountPrecision: int64(3),
			PricePrecision:  int64(6),
		}, odxm)
	})

	t.Run("BTCUSDT", func(t *testing.T) {
		bm := &BinanceMarket{
			Symbol:         "BTCUSDT",
			BaseUnit:       "BTC",
			QuoteUnit:      "USDT",
			QuotePrecision: decimal.RequireFromString("8"),
			Filters: []Filter{
				{
					Type:     "PRICE_FILTER",
					MinPrice: decimal.RequireFromString("0.01000000"),
					MaxPrice: decimal.RequireFromString("1000000.00000000"),
					TickSize: decimal.RequireFromString("0.01000000"),
				},
				{
					Type:        "LOT_SIZE",
					MinQuantity: decimal.RequireFromString("0.00002000"),
				},
				{
					Type:        "MIN_NOTIONAL",
					MinNotional: decimal.RequireFromString("10.00000000"),
				},
			},
		}

		odxm, err := bm.ToOpendaxMarket(decimal.NewFromFloat32(float32(0.0002))) // BTC at 50k
		require.NoError(t, err)

		expectedBM := &opendax.OpendaxMarket{
			Symbol:          "btcusdt",
			Name:            "BTC/USDT",
			BaseUnit:        "btc",
			QuoteUnit:       "usdt",
			MinPrice:        decimal.RequireFromString("0.01000000"),
			MaxPrice:        decimal.RequireFromString("0.00"),
			MinAmount:       decimal.RequireFromString("0.00020"),
			AmountPrecision: int64(5),
			PricePrecision:  int64(2),
		}

		fmt.Printf("EXPECTED BINANCE MARKET: %v\n", odxm)
		fmt.Printf("ACTUAL BINANCE MARKET:   %v\n", expectedBM)

		assert.DeepEqual(t, expectedBM, odxm)

		odxm, err = bm.ToOpendaxMarket(decimal.NewFromFloat32(float32(0.000005))) // BTC at 2M
		require.NoError(t, err)

		expectedMarket := &opendax.OpendaxMarket{
			Symbol:          "btcusdt",
			Name:            "BTC/USDT",
			BaseUnit:        "btc",
			QuoteUnit:       "usdt",
			MinPrice:        decimal.RequireFromString("0.01"),
			MaxPrice:        decimal.RequireFromString("0"),
			MinAmount:       decimal.RequireFromString("0.00002"),
			AmountPrecision: int64(5),
			PricePrecision:  int64(2),
		}

		fmt.Printf("EXPECTED MARKET: %v\n", expectedMarket)
		fmt.Printf("ACTUAL MARKET:   %v\n", odxm)

		assert.DeepEqual(t, expectedMarket, odxm)

	})

	/*
		t.Run("XRPUSDT", func(t *testing.T) {
			bm := &BinanceMarket{
				Symbol:         "XRPUSDT",
				BaseUnit:       "XRP",
				QuoteUnit:      "USDT",
				QuotePrecision: decimal.RequireFromString("8"),
				Filters: []Filter{
					{
						Type:     "PRICE_FILTER",
						MinPrice: decimal.RequireFromString("0.00010000"),
						MaxPrice: decimal.RequireFromString("10000.00000000"),
						TickSize: decimal.RequireFromString("0.00010000"),
					},
					{
						Type:        "LOT_SIZE",
						MinQuantity: decimal.RequireFromString("1.00000000"),
					},
					{
						Type:        "MIN_NOTIONAL",
						MinNotional: decimal.RequireFromString("10.00000000"),
					},
				},
			}

			odxm, err := bm.ToOpendaxMarket(decimal.NewFromFloat32(9.42)) // XRP at 1.1143
			require.NoError(t, err)

			assert.Equal(t, &opendax.OpendaxMarket{
				Symbol:          "xrpusdt",
				Name:            "XRP/USDT",
				BaseUnit:        "xrp",
				QuoteUnit:       "usdt",
				MinPrice:        decimal.RequireFromString("0.00010000"),
				MaxPrice:        decimal.RequireFromString("0.00"),
				MinAmount:       decimal.RequireFromString("10"),
				AmountPrecision: int64(0),
				PricePrecision:  int64(4),
			}, odxm)
		})
	*/
}
