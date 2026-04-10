package main

import (
	"fmt"
	"log"

	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	ws := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/spot", func(message string) error {
		fmt.Println("Received:", message)
		return nil
	})
	if err := ws.Connect(); err != nil {
		log.Fatal(err)
	}
	if _, err := ws.SendSubscription([]string{"orderbook.1.BTCUSDT", "orderbook.1.ETHUSDT"}); err != nil {
		log.Fatal(err)
	}
	select {}
}
