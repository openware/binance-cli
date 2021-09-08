package opendax

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func NewOpendaxClient(platformUrl string) *OpendaxClient {
	return &OpendaxClient{
		platformUrl: platformUrl,
	}
}

func (oc *OpendaxClient) Authorize(apiKey, secretKey string) {
	oc.apiKey = apiKey
	oc.secretKey = secretKey
}

func (oc *OpendaxClient) FetchOpendaxCurrencies() (OpendaxCurrencies, error) {
	currencies := OpendaxCurrencies{}
	_, _, _, err := oc.opendaxApiCall(currenciesEndpoint, &currencies)
	return currencies, err
}

func (oc *OpendaxClient) FetchOpendaxMarkets() (OpendaxMarkets, error) {
	markets := OpendaxMarkets{}
	_, _, _, err := oc.opendaxApiCall(marketsEndpoint, &markets)
	return markets, err
}

func (oc *OpendaxClient) UpdateOpendaxMarket(request UpdateMarketRequest) (OpendaxMarket, error) {
	body, err := request.Encode()
	if err != nil {
		panic(err)
	}

	market := OpendaxMarket{}
	_, _, _, err = oc.opendaxPostApiCall(adminMarketsUpdateEndpoint, body, &market)
	return market, err
}

func (oc *OpendaxClient) UpdateOpendaxSecret(request UpdateSecretRequest) error {
	body, err := request.Encode()
	if err != nil {
		panic(err)
	}

	_, _, _, err = oc.opendaxPostApiCall(adminFinexSecretUpdateEndpoint, body, nil)
	return err
}

func (oc *OpendaxClient) SignRequest(request *http.Request) {
	nonce := fmt.Sprintf("%d", time.Now().Unix()*1000)

	data := strings.Join([]string{nonce, oc.apiKey}, "")

	h := hmac.New(sha256.New, []byte(oc.secretKey))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))

	request.Header.Set("X-Auth-Apikey", oc.apiKey)
	request.Header.Set("X-Auth-Nonce", nonce)
	request.Header.Set("X-Auth-Signature", sha)
}
