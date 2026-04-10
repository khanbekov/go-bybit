package main

import (
	"fmt"

	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	ws := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/linear", func(message string) error {
		fmt.Println("Received:", message)
		return nil
	})
	_, _ = ws.Connect().SubscribeKline(bybit.KlineInterval1m, "BTCUSDT", "ETHUSDT")
	select {}
}
