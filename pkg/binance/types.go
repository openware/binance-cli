package binance

import (
	"fmt"
	"strings"

	"github.com/openware/binance-cli/pkg/helpers"
	"github.com/openware/binance-cli/pkg/opendax"
	"github.com/shopspring/decimal"
)

type BinanceClient struct {
	apiKey string
	secret string
	url    string
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
	Symbol         string          `json:"symbol"`
	BaseUnit       string          `json:"baseAsset"`
	QuoteUnit      string          `json:"quoteAsset"`
	QuotePrecision decimal.Decimal `json:"quotePrecision"`
	Filters        []Filter        `json:"filters"`
}

type Filter struct {
	Type        string          `json:"filterType"`
	MinPrice    decimal.Decimal `json:"minPrice,omitempty"`
	MaxPrice    decimal.Decimal `json:"maxPrice,omitempty"`
	TickSize    decimal.Decimal `json:"tickSize,omitempty"`
	MinQuantity decimal.Decimal `json:"minQty,omitempty"`
	MinNotional decimal.Decimal `json:"minNotional,omitempty"`
}

func (m *BinanceMarket) GetFilter(filterType string) (*Filter, error) {
	for _, f := range m.Filters {
		if f.Type == filterType {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("Filter %s not found", filterType)
}

type BinanceCurrencies []*BinanceCurrency

type BinanceCurrency struct {
	Code     string    `json:"coin"`
	Networks []Network `json:"networkList"`
}

type Network struct {
	Name        string          `json:"network"`
	WithdrawFee decimal.Decimal `json:"withdrawFee"`
	WithdrawMin decimal.Decimal `json:"withdrawMin"`
}

type BinanceTickerPrice struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
}

func (m *BinanceMarket) ToOpendaxMarket(minAmount decimal.Decimal) (*opendax.OpendaxMarket, error) {
	priceFilter, err := m.GetFilter("PRICE_FILTER")
	if err != nil {
		return nil, err
	}

	quantityFilter, err := m.GetFilter("LOT_SIZE")
	if err != nil {
		return nil, err
	}

	tickPrecision := decimal.NewFromInt(helpers.ValuePrecision(priceFilter.TickSize))

	pricePrecision := decimal.Min(tickPrecision, m.QuotePrecision)
	amountPrecision := helpers.ValuePrecision(quantityFilter.MinQuantity)

	minAmount = minAmount.Round(int32(amountPrecision))

	if minAmount.LessThan(quantityFilter.MinQuantity) {
		minAmount = quantityFilter.MinQuantity
	}

	return &opendax.OpendaxMarket{
		Symbol:          strings.ToLower(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "")),
		Name:            strings.ToUpper(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "/")),
		BaseUnit:        strings.ToLower(m.BaseUnit),
		QuoteUnit:       strings.ToLower(m.QuoteUnit),
		MinPrice:        priceFilter.MinPrice,
		MaxPrice:        decimal.Zero,
		MinAmount:       minAmount,
		AmountPrecision: amountPrecision,
		PricePrecision:  pricePrecision.IntPart(),
	}, nil
}

func (m *BinanceMarket) OpendaxMarketName() string {
	return strings.ToUpper(strings.Join([]string{m.BaseUnit, m.QuoteUnit}, "/"))
}

func (m BinanceMarket) CalculateMinAmount(price decimal.Decimal) decimal.Decimal {
	notionalFilter, err := m.GetFilter("MIN_NOTIONAL")
	if err != nil {
		return decimal.Zero
	}

	// Return 105% of min amount to be sure it covers the min notional
	return decimal.RequireFromString("1.05").Mul(notionalFilter.MinNotional).Div(price)
}
