package main

import (
	"fmt"
	"log"

	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	ws := bybit.NewBybitPublicWebSocket(bybit.SPREAD_MAINNET, func(message string) error {
		fmt.Println("Received:", message)
		return nil
	})
	if err := ws.Connect(); err != nil {
		log.Fatal(err)
	}
	if _, err := ws.SendSubscription([]string{"orderbook.25.SOLUSDT_SOL/USDT"}); err != nil {
		log.Fatal(err)
	}
	select {}
}
