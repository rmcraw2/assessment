package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rmcraw2/assessment/internal/client/coinbase"
)

type coinbaseInterface interface {
	GetExchangeRates() (*coinbase.ExchangeRates, error)
}

func main() {
	cbClient := coinbase.CoinbaseClient{}

	line1, line2, err := MainBody(cbClient)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*line1)
	fmt.Println(*line2)
}

func MainBody(cbClient coinbaseInterface) (*string, *string, error) {
	if len(os.Args) <= 3 {
		return nil, nil, errors.New("not enough inputs")
	}

	usdAmount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		return nil, nil, err
	}

	coin1 := strings.ToUpper(os.Args[2])
	coin2 := strings.ToUpper(os.Args[3])

	exchangeRates, err := cbClient.GetExchangeRates()
	if err != nil {
		return nil, nil, err
	}

	amount1 := roundFloat2(usdAmount * 0.70)
	amount2 := roundFloat2(usdAmount - amount1)

	coinAmount1, err := exchangeCoin(amount1, coin1, exchangeRates)
	if err != nil {
		return nil, nil, err
	}

	coinAmount2, err := exchangeCoin(amount2, coin2, exchangeRates)
	if err != nil {
		return nil, nil, err
	}

	line1 := fmt.Sprintf("$%.2f => %.4f %v", amount1, coinAmount1, coin1)
	line2 := fmt.Sprintf("$%.2f => %.4f %v", amount2, coinAmount2, coin2)

	return &line1, &line2, nil
}

func roundFloat2(input float64) float64 {
	floatString := fmt.Sprintf("%.2f", input)
	output, _ := strconv.ParseFloat(floatString, 64)
	return output
}

func exchangeCoin(amount float64, coin string, exchangeRates *coinbase.ExchangeRates) (float64, error) {
	rate, ok := exchangeRates.Rates[coin]
	if !ok {
		return 0, errors.New("unknown coin")
	}

	rateFloat, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return 0, err
	}

	return amount * rateFloat, nil
}
