package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	GetCoinInfo()
}

func GetCoinInfo() {
	client := bybit.NewBybitHttpClient("YOUR_API_KEY", "YOUR_API_SECRET", bybit.WithBaseURL(bybit.TESTNET))
	Params := map[string]interface{}{"coin": "USDT"}
	assetResult, err := client.NewUtaBybitServiceWithParams(Params).GetCoinInfo(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(assetResult))
}
