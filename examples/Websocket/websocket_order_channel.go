package main

import (
	"fmt"
	"log"

	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	ws := bybit.NewBybitPrivateWebSocket("wss://stream-testnet.bybit.com/v5/private", "YOUR_API_KEY", "YOUR_API_SECRET", func(message string) error {
		fmt.Println("Received:", message)
		return nil
	}, bybit.WithPingInterval(2))
	if err := ws.Connect(); err != nil {
		log.Fatal(err)
	}
	if _, err := ws.SendSubscription([]string{"order", "position", "wallet"}); err != nil {
		log.Fatal(err)
	}
	select {}
}
