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
}

func (cb CoinbaseClient) GetExchangeRates() (*ExchangeRates, error) {
	url := fmt.Sprintf("https://api.coinbase.com/v2/exchange-rates?currency=USD")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	exchangeRatesData := &ExchangeRatesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(exchangeRatesData); err != nil {
		return nil, err
	}

	return &exchangeRatesData.Data, nil
}
