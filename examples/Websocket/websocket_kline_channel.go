package main

import (
	"fmt"
	"log"

	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	ws := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/linear", func(message string) error {
		fmt.Println("Received:", message)
		return nil
	})
	if err := ws.Connect(); err != nil {
		log.Fatal(err)
	}
	if _, err := ws.SubscribeKline(bybit.KlineInterval1m, "BTCUSDT", "ETHUSDT"); err != nil {
		log.Fatal(err)
	}
	select {}
}
