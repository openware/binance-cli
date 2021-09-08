package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	binanceBaseUrl          = "https://api.binance.com"
	coinsInfoEndpoint       = "/sapi/v1/capital/config/getall"
	exchangeInfoEndpoint    = "/api/v3/exchangeInfo"
	tickerPriceInfoEndpoint = "/api/v3/ticker/price"
	NotFoundError           = "404 Record Not Found"
	ServiceUnavailableError = "503 Service Unavailable"
	HttpTransportError      = "HTTP Transport Error"
)

var SignedEndpoints = map[string]struct{}{
	coinsInfoEndpoint: {},
}

func (bc *BinanceClient) apiCall(endpoint string, receiver interface{}) (interface{}, error) {
	uri := binanceBaseUrl + endpoint
	req, _ := http.NewRequest("GET", uri, nil)

	fmt.Printf("Calling %s\n", uri)

	req.Header.Add("X-MBX-APIKEY", bc.apiKey)

	if _, ok := SignedEndpoints[endpoint]; ok {
		q, err := bc.mandatoryBinanceParameters()
		if err != nil {
			return receiver, err
		}
		req.URL.RawQuery = q
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return receiver, fmt.Errorf(NotFoundError)
	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusBadGateway {
		return receiver, fmt.Errorf(ServiceUnavailableError)
	}

	err = json.NewDecoder(resp.Body).Decode(receiver)

	re, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(re))
	return receiver, err
}

func (bc *BinanceClient) mandatoryBinanceParameters() (string, error) {
	timestampParam := fmt.Sprintf("timestamp=%v", currentTimestamp())
	mac := hmac.New(sha256.New, []byte(bc.secret))
	_, err := mac.Write([]byte(timestampParam))
	if err != nil {
		return "", err
	}

	signature := fmt.Sprintf("%x", mac.Sum(nil))
	return fmt.Sprintf("timestamp=%v&signature=%s", currentTimestamp(), signature), nil
}

func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
