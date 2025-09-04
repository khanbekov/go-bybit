package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	GetMarketKline()
}

func GetMarketKline() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.TESTNET))
	Params := map[string]interface{}{"category": "spot", "symbol": "BTCUSDT", "interval": "1"}
	serverResult, err := client.NewUtaBybitServiceWithParams(Params).GetMarketKline(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}
