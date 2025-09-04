package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	GetDeliveryPrice()
}

func GetDeliveryPrice() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.TESTNET))
	Params := map[string]interface{}{"category": "linear", "symbol": "BTCUSDT"}
	serverResult, err := client.NewUtaBybitServiceWithParams(Params).GetDeliveryPrice(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}
