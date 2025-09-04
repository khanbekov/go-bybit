package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	GetIndexPriceKline()
}

func GetIndexPriceKline() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.TESTNET))
	Params := map[string]interface{}{"category": "linear", "symbol": "BTCUSDTT", "interval": "1"}
	serverResult, err := client.NewUtaBybitServiceWithParams(Params).GetIndexPriceKline(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}
