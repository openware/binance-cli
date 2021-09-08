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
	MinPrice    json.Number `json:"minPrice,omitempty"`
	MaxPrice    json.Number `json:"maxPrice,omitempty"`
	TickSize    json.Number `json:"tickSize,omitempty"`
	MinQuantity json.Number `json:"minQty,omitempty"`
	MinNotional json.Number `json:"minNotional,omitempty"`
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

type BinanceTickerPrice struct {
	Symbol string      `json:"symbol"`
	Price  json.Number `json:"price"`
}

func (m BinanceMarket) ToOpendaxMarket(minAmountFloat float64) opendax.OpendaxMarket {
	priceFilter := m.GetFilter("PRICE_FILTER")
	quantityFilter := m.GetFilter("LOT_SIZE")

	tickPrecision := helpers.ValuePrecision(priceFilter.TickSize)
	quotePrecision, _ := m.QuotePrecision.Float64()

	pricePresion := math.Min(float64(tickPrecision), quotePrecision)
	amountPrecision := helpers.ValuePrecision(quantityFilter.MinQuantity)

	minAmount := json.Number(fmt.Sprintf("%0."+fmt.Sprint(amountPrecision)+"f", minAmountFloat))

	return opendax.OpendaxMarket{
		Symbol:    strings.ToLower(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "")),
		Name:      strings.ToUpper(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "/")),
		BaseUnit:  strings.ToLower(m.BaseUnit),
		QuoteUnit: strings.ToLower(m.QuoteUnit),
		MinPrice:  priceFilter.MinPrice,
		// Commented to not limit MaxPrice
		//MaxPrice:        priceFilter.MaxPrice,
		MaxPrice:        json.Number(fmt.Sprintf(`%.2f`, 0.0)),
		MinAmount:       minAmount,
		AmountPrecision: amountPrecision,
		PricePrecision:  int(pricePresion),
	}
}

func (m BinanceMarket) OpendaxMarketName() string {
	return strings.ToUpper(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "/"))
}

func (m BinanceMarket) CalculateMinAmount(price json.Number) (float64, error) {
	notionalFilter := m.GetFilter("MIN_NOTIONAL")

	minNotionalFloat, err := notionalFilter.MinNotional.Float64()
	if err != nil {
		return 0, err
	}

	priceFloat, err := price.Float64()
	if err != nil {
		return 0, err
	}

	// Return 105% of min amount to be sure it covers the min notional
	minAmount := 1.05 * minNotionalFloat / priceFloat

	return minAmount, nil
}
