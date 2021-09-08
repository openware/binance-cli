package opendax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	adminMarketsUpdateEndpoint     = "/api/v2/peatio/admin/markets/update"
	adminFinexSecretUpdateEndpoint = "/api/v2/sonic/admin/finex/secret"
	marketsEndpoint                = "/api/v2/peatio/public/markets"
	currenciesEndpoint             = "/api/v2/peatio/public/currencies"
	NotFoundError                  = "404 Record Not Found"
	NotAuthorizedError             = "401 Not Authorized"
	ServiceUnavailableError        = "503 Service Unavailable"
	HttpTransportError             = "HTTP Transport Error"
)

func (oc *OpendaxClient) opendaxApiCall(endpoint string, receiver interface{}) (interface{}, http.Header, int, error) {
	uri := oc.platformUrl + endpoint
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return receiver, http.Header{}, 0, fmt.Errorf(HttpTransportError)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return receiver, resp.Header, resp.StatusCode, fmt.Errorf(NotFoundError)
	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusBadGateway {
		return receiver, resp.Header, resp.StatusCode, fmt.Errorf(ServiceUnavailableError)
	}

	err = json.NewDecoder(resp.Body).Decode(receiver)
	return receiver, resp.Header, resp.StatusCode, err
}

func (oc *OpendaxClient) opendaxPostApiCall(endpoint string, body []byte, receiver interface{}) (interface{}, http.Header, int, error) {
	uri := oc.platformUrl + endpoint

	// TODO: Refactor to pass method into opendaxPostApiCall
	method := "POST"
	if endpoint == adminFinexSecretUpdateEndpoint {
		method = "PUT"
	}

	req, err := http.NewRequest(method, uri, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	oc.SignRequest(req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return receiver, http.Header{}, 0, fmt.Errorf(HttpTransportError)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		return receiver, resp.Header, resp.StatusCode, fmt.Errorf(NotFoundError)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return receiver, resp.Header, resp.StatusCode, fmt.Errorf(NotAuthorizedError)
	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusBadGateway {
		return receiver, resp.Header, resp.StatusCode, fmt.Errorf(ServiceUnavailableError)
	}

	err = json.NewDecoder(resp.Body).Decode(receiver)
	return receiver, resp.Header, resp.StatusCode, err
}
