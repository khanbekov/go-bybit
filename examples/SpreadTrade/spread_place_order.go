package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	SpreadTradeOrderMap()
}
func SpreadTradeOrderMap() {
	client := bybit.NewBybitHttpClient("YOUR_API_KEY", "YOUR_API_SECRET", bybit.WithBaseURL(bybit.TESTNET))
	Params := map[string]interface{}{
		"symbol":      "SOLUSDT_SOL/USDT",
		"side":        "Buy",
		"orderType":   "Limit",
		"qty":         "0.1",
		"price":       "21",
		"orderLinkId": "1744072052193428479",
		"timeInForce": "PostOnly",
	}
	orderResult, err := client.NewUtaBybitServiceWithParams(Params).PlaceSpreadTradeOrder(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(orderResult))
}
