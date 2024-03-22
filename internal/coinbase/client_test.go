package coinbase_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rmcraw2/assessment/internal/coinbase"
	"github.com/stretchr/testify/suite"
)

type coinbaseSuite struct {
	suite.Suite
}

func TestCoinbaseSuite(t *testing.T) {
	suite.Run(t, new(coinbaseSuite))
}

func (suite *coinbaseSuite) TestGetExchangeRates() {
	testCases := map[string]struct {
		testServer  func() *httptest.Server
		testAsserts func(*coinbase.ExchangeRates, error)
	}{
		"Test success": {
			testServer: func() *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(getExchangeRatesResponse())
				}))

				return server
			},
			testAsserts: func(exchangeRates *coinbase.ExchangeRates, err error) {
				suite.Nil(err)
				suite.NotNil(exchangeRates)
			},
		},
		"Test failure": {
			testServer: func() *httptest.Server {
				return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			},
			testAsserts: func(exchangeRates *coinbase.ExchangeRates, err error) {
				suite.Equal(`Get "/v2/exchange-rates?currency=USD": unsupported protocol scheme ""`, err.Error())
				suite.Nil(exchangeRates)
			},
		},
	}

	for tn, tc := range testCases {
		suite.Run(tn, func() {
			server := tc.testServer()

			cbClient := coinbase.NewCoinbaseClient()
			cbClient.URL = server.URL

			exchangeRates, err := cbClient.GetExchangeRates()

			tc.testAsserts(exchangeRates, err)
		})
	}
}

func getExchangeRatesResponse() []byte {
	exchangeRatesResponse := coinbase.ExchangeRatesResponse{
		Data: coinbase.ExchangeRates{
			Currency: "USD",
			Rates: map[string]string{
				"BTC": "0.0000158352814028",
				"ETH": "0.0003074515494789",
			},
		},
	}

	resp, _ := json.Marshal(exchangeRatesResponse)
	return resp
}
