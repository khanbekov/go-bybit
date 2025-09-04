package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	GetHistoryVolatility()
}

func GetHistoryVolatility() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.TESTNET), bybit.WithDebug(true))
	Params := map[string]interface{}{"category": "option", "baseCoin": "BTC"}
	serverResult, err := client.NewUtaBybitServiceWithParams(Params).GetHistoryVolatility(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}
