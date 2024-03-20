package main_test

import (
	"os"
	"testing"

	"github.com/rmcraw2/assessment/cmd/main"
	"github.com/rmcraw2/assessment/internal/coinbase"
	"github.com/stretchr/testify/suite"
)

type mainSuite struct {
	suite.Suite
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(mainSuite))
}

func (suite *mainSuite) TestMainBody() {

	testCases := map[string]struct {
		testInputs  func() []string
		cbClient    func() main.CoinbaseInterface
		testAsserts func(*string, *string, error)
	}{
		"$100 to BTC and ETH": {
			testInputs: func() []string {
				return []string{"main", "100", "BTC", "ETH"}
			},
			cbClient: func() main.CoinbaseInterface {
				return mockCbClient{mockGetExchangeRates: func() (*coinbase.ExchangeRates, error) {
					return &coinbase.ExchangeRates{
						Currency: "USD",
						Rates: map[string]string{
							"BTC": "0.0000158352814028",
							"ETH": "0.0003074515494789",
						},
					}, nil
				}}
			},
			testAsserts: func(line1 *string, line2 *string, err error) {
				suite.Equal("$70.00 => 0.0011 BTC", *line1)
				suite.Equal("$30.00 => 0.0092 ETH", *line2)
				suite.Nil(err)
			},
		},
		"$0.01 to BTC and ETH": {
			testInputs: func() []string {
				return []string{"main", "0.01", "BTC", "ETH"}
			},
			cbClient: func() main.CoinbaseInterface {
				return mockCbClient{mockGetExchangeRates: func() (*coinbase.ExchangeRates, error) {
					return &coinbase.ExchangeRates{
						Currency: "USD",
						Rates: map[string]string{
							"BTC": "0.0000158352814028",
							"ETH": "0.0003074515494789",
						},
					}, nil
				}}
			},
			testAsserts: func(line1 *string, line2 *string, err error) {
				suite.Equal("$0.01 => 0.0000 BTC", *line1)
				suite.Equal("$0.00 => 0.0000 ETH", *line2)
				suite.Nil(err)
			},
		},
		"$0.001 to BTC and ETH": {
			testInputs: func() []string {
				return []string{"main", "0.001", "BTC", "ETH"}
			},
			cbClient: func() main.CoinbaseInterface {
				return mockCbClient{mockGetExchangeRates: func() (*coinbase.ExchangeRates, error) {
					return &coinbase.ExchangeRates{
						Currency: "USD",
						Rates: map[string]string{
							"BTC": "0.0000158352814028",
							"ETH": "0.0003074515494789",
						},
					}, nil
				}}
			},
			testAsserts: func(line1 *string, line2 *string, err error) {
				suite.Equal("$0.00 => 0.0000 BTC", *line1)
				suite.Equal("$0.00 => 0.0000 ETH", *line2)
				suite.Nil(err)
			},
		},
		"$1000000 to BTC and ETH": {
			testInputs: func() []string {
				return []string{"main", "1000000", "BTC", "ETH"}
			},
			cbClient: func() main.CoinbaseInterface {
				return mockCbClient{mockGetExchangeRates: func() (*coinbase.ExchangeRates, error) {
					return &coinbase.ExchangeRates{
						Currency: "USD",
						Rates: map[string]string{
							"BTC": "0.0000158352814028",
							"ETH": "0.0003074515494789",
						},
					}, nil
				}}
			},
			testAsserts: func(line1 *string, line2 *string, err error) {
				suite.Equal("$700000.00 => 11.0847 BTC", *line1)
				suite.Equal("$300000.00 => 92.2355 ETH", *line2)
				suite.Nil(err)
			},
		},
		"Not enough args error": {
			testInputs: func() []string {
				return []string{"main"}
			},
			cbClient: func() main.CoinbaseInterface {
				return mockCbClient{mockGetExchangeRates: func() (*coinbase.ExchangeRates, error) {
					return &coinbase.ExchangeRates{
						Currency: "USD",
						Rates:    map[string]string{},
					}, nil
				}}
			},
			testAsserts: func(line1 *string, line2 *string, err error) {
				suite.Nil(line1)
				suite.Nil(line2)
				suite.Equal("not enough inputs", err.Error())
			},
		},
		"unknown coin": {
			testInputs: func() []string {
				return []string{"main", "100", "X", "Z"}
			},
			cbClient: func() main.CoinbaseInterface {
				return mockCbClient{mockGetExchangeRates: func() (*coinbase.ExchangeRates, error) {
					return &coinbase.ExchangeRates{
						Currency: "USD",
						Rates:    map[string]string{},
					}, nil
				}}
			},
			testAsserts: func(line1 *string, line2 *string, err error) {
				suite.Nil(line1)
				suite.Nil(line2)
				suite.Equal("unknown coin", err.Error())
			},
		},
	}

	for tn, tc := range testCases {
		suite.Run(tn, func() {
			args := tc.testInputs()
			os.Args = args

			line1, line2, err := main.MainBody(tc.cbClient())

			tc.testAsserts(line1, line2, err)
		})
	}

}

type mockCbClient struct {
	mockGetExchangeRates func() (*coinbase.ExchangeRates, error)
}

func (mcb mockCbClient) GetExchangeRates() (*coinbase.ExchangeRates, error) {
	return mcb.mockGetExchangeRates()
}
