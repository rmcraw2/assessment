package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rmcraw2/assessment/internal/client/coinbase"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter USD for conversion: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	exchangeRates, err := coinbase.GetExchangeRates("USD")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*exchangeRates)
}
