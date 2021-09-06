package binance

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/openware/binance-cli/pkg/helpers"
	"github.com/openware/binance-cli/pkg/opendax"
)

type BinanceClient struct {
	apiKey string
	secret string
}

type BinanceExchangeInfo struct {
	Symbols        []BinanceMarket `json:"symbols"`
	MarketRegistry map[string]BinanceMarket
}

func (info *BinanceExchangeInfo) FillRegistry() {
	info.MarketRegistry = make(map[string]BinanceMarket)
	for _, m := range info.Symbols {
		info.MarketRegistry[m.Symbol] = m
	}
}

type BinanceMarket struct {
	Symbol         string      `json:"symbol"`
	BaseUnit       string      `json:"baseAsset"`
	QuoteUnit      string      `json:"quoteAsset"`
	QuotePrecision json.Number `json:"quotePrecision"`
	Filters        []Filter    `json:"filters"`
}

type Filter struct {
	Type        string      `json:"filterType"`
	MinPrice    json.Number `json:"minPrice"`
	MaxPrice    json.Number `json:"maxPrice"`
	TickSize    json.Number `json:"tickSize"`
	MinQuantity json.Number `json:"minQty"`
}

func (m BinanceMarket) GetFilter(filter string) Filter {
	for _, f := range m.Filters {
		if f.Type == filter {
			return f
		}
	}

	return Filter{}
}

type BinanceCurrencies []*BinanceCurrency

type BinanceCurrency struct {
	Code     string    `json:"coin"`
	Networks []Network `json:"networkList"`
}

type Network struct {
	Name        string      `json:"network"`
	WithdrawFee json.Number `json:"withdrawFee"`
	WithdrawMin json.Number `json:"withdrawMin"`
}

func (m BinanceMarket) ToOpendaxMarket() opendax.OpendaxMarket {
	priceFilter := m.GetFilter("PRICE_FILTER")
	quantityFilter := m.GetFilter("LOT_SIZE")

	tickPrecision := helpers.ValuePrecision(priceFilter.TickSize)
	quotePrecision, _ := m.QuotePrecision.Float64()

	pricePresion := math.Min(float64(tickPrecision), quotePrecision)

	return opendax.OpendaxMarket{
		Symbol:    strings.ToLower(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "")),
		Name:      strings.ToUpper(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "/")),
		BaseUnit:  strings.ToLower(m.BaseUnit),
		QuoteUnit: strings.ToLower(m.QuoteUnit),
		MinPrice:  priceFilter.MinPrice,
		// Commented to not limit MaxPrice
		//MaxPrice:        priceFilter.MaxPrice,
		MaxPrice:        json.Number(fmt.Sprintf(`%.2f`, 0.0)),
		MinAmount:       quantityFilter.MinQuantity,
		AmountPrecision: helpers.ValuePrecision(quantityFilter.MinQuantity),
		PricePrecision:  int(pricePresion),
	}
}

func (m BinanceMarket) OpendaxMarketName() string {
	return strings.ToUpper(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "/"))
}
