package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	GetLongShortRatio()
}

func GetLongShortRatio() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.TESTNET))
	Params := map[string]interface{}{"category": "linear", "symbol": "BTCUSDT", "period": "5min"}
	serverResult, err := client.NewUtaBybitServiceWithParams(Params).GetLongShortRatio(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}
