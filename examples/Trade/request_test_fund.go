package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	RequestTestFund()
}
func RequestTestFund() {
	client := bybit.NewBybitHttpClient("xxxx", "xxx", bybit.WithBaseURL(bybit.DEMO_ENV))
	Params := map[string]interface{}{
		"adjustType": 0,
		"utaDemoApplyMoney": []map[string]interface{}{
			{
				"coin":      "USDT",
				"amountStr": "109",
			},
			{
				"coin":      "ETH",
				"amountStr": "1",
			},
		},
	}
	serverResult, err := client.NewUtaBybitServiceWithParams(Params).RequestTestFund(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(serverResult))
}
