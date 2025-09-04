package main

import (
	"context"
	"fmt"
	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	PlaceOrder()
}

func PlaceOrder() {
	client := bybit.NewBybitHttpClient("YOUR_API_KEY", "YOUR_API_SECRET", bybit.WithBaseURL(bybit.MAINNET), bybit.WithDebug(true))
	orderResult, err := client.NewPlaceOrderService("linear", "RUNEUSDT", "Buy", "LIMIT", "15").Price("1.1").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(orderResult))
}
