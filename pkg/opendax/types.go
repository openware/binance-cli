package opendax

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type OpendaxClient struct {
	platformUrl string
	apiKey      string
	secretKey   string
}

type OpendaxCurrencies []*OpendaxCurrency

type OpendaxCurrency struct {
	Code              string          `json:"id"`
	WithdrawFee       decimal.Decimal `json:"withdraw_fee"`
	MinWithdrawAmount decimal.Decimal `json:"min_withdraw_amount"`
}

func (c *OpendaxCurrency) ToBinanceCoinName() string {
	return strings.ToUpper(c.Code)
}

type OpendaxMarkets []OpendaxMarket

type OpendaxMarket struct {
	Symbol          string          `json:"symbol"`
	Name            string          `json:"name"`
	BaseUnit        string          `json:"base_unit"`
	QuoteUnit       string          `json:"quote_unit"`
	MinPrice        decimal.Decimal `json:"min_price"`
	MaxPrice        decimal.Decimal `json:"max_price"`
	MinAmount       decimal.Decimal `json:"min_amount"`
	AmountPrecision int64           `json:"amount_precision"`
	PricePrecision  int64           `json:"price_precision"`
}

func (om OpendaxMarket) ToBinanceMarketName() string {
	return strings.ToUpper(strings.Join([]string{om.BaseUnit, om.QuoteUnit}, ""))
}

func (om OpendaxMarket) Print() {
	fmt.Println("- 	Symbol:", om.Symbol)
	fmt.Println("	Name:", om.Name)
	fmt.Println("	BaseUnit:", om.BaseUnit)
	fmt.Println("	QuoteUnit:", om.QuoteUnit)
	fmt.Println("	MinPrice:", om.MinPrice)
	fmt.Println("	MaxPrice:", om.MaxPrice)
	fmt.Println("	MinAmount:", om.MinAmount)
	fmt.Println("	AmountPrecision:", om.AmountPrecision)
	fmt.Println("	PricePrecision:", om.PricePrecision)
	fmt.Println("")
}

type Request interface {
	Encode() ([]byte, error)
}

type UpdateMarketRequest struct {
	Symbol          string          `json:"symbol"`
	MinPrice        decimal.Decimal `json:"min_price"`
	MaxPrice        decimal.Decimal `json:"max_price"`
	MinAmount       decimal.Decimal `json:"min_amount"`
	AmountPrecision int64           `json:"amount_precision"`
	PricePrecision  int64           `json:"price_precision"`
}

func (r *UpdateMarketRequest) Encode() ([]byte, error) {
	return json.Marshal(r)
}

func CompareOpendaxMarkets(firstMarket, secondMarket *OpendaxMarket) bool {
	return firstMarket.AmountPrecision == secondMarket.AmountPrecision && firstMarket.PricePrecision == secondMarket.PricePrecision && firstMarket.MinAmount.Equals(secondMarket.MinAmount)
}

// UpdateSecretRequest represents params for a Sonic secret update request
type UpdateSecretRequest struct {
	Key   string `json:"key"`
	Scope string `json:"scope"`
	Value string `json:"value"`
}

func (r *UpdateSecretRequest) Encode() ([]byte, error) {
	return json.Marshal(r)
}
