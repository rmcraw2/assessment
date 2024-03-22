package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// https://docs.cloud.coinbase.com/sign-in-with-coinbase/docs/api-exchange-rates

type ExchangeRatesResponse struct {
	Data ExchangeRates `json:"data"`
}

type ExchangeRates struct {
	Currency string            `json:"currency"`
	Rates    map[string]string `json:"rates"`
}

type CoinbaseClient struct {
	URL string
}

func NewCoinbaseClient() *CoinbaseClient {
	return &CoinbaseClient{
		URL: "https://api.coinbase.com",
	}
}

func (cb CoinbaseClient) GetExchangeRates() (*ExchangeRates, error) {
	url := fmt.Sprintf("%v/v2/exchange-rates?currency=USD", cb.URL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ExchangeRatesResponse := &ExchangeRatesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(ExchangeRatesResponse); err != nil {
		return nil, err
	}

	return &ExchangeRatesResponse.Data, nil
}
