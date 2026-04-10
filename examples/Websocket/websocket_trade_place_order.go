package main

import (
	"fmt"
	"log"
	"time"

	bybit "github.com/khanbekov/go-bybit"
)

func main() {
	ws := bybit.NewBybitPrivateWebSocket(bybit.WEBSOCKET_TRADE_TESTNET, "YOUR_API_KEY", "YOUR_API_SECRET", func(message string) error {
		fmt.Println("Received:", message)
		return nil
	}, bybit.WithPingInterval(10))
	if err := ws.Connect(); err != nil {
		log.Fatal(err)
	}
	if _, err := ws.SendTradeRequest(map[string]interface{}{
		"reqId": "test-005",
		"header": map[string]string{
			"X-BAPI-TIMESTAMP":   fmt.Sprintf("%d", time.Now().UnixMilli()),
			"X-BAPI-RECV-WINDOW": "8000",
			"Referer":            "bot-001",
		},
		"op": "order.create",
		"args": []interface{}{
			map[string]interface{}{
				"symbol":      "ETHUSDT",
				"side":        "Buy",
				"orderType":   "Limit",
				"qty":         "0.2",
				"price":       "2800",
				"category":    "linear",
				"timeInForce": "PostOnly",
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
	select {}
}
