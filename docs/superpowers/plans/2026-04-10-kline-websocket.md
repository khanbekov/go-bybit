# Kline Websocket Helper Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a typed `SubscribeKline` helper, interval constants, a parsed event model, and an example to cover the only missing endpoint in the audit — real-time kline websocket.

**Architecture:** Keep the generic `WebSocket.SendSubscription` unchanged. Add a thin typed layer on top: `KlineInterval` string type with constants, an `IsValid()` check, a `SubscribeKline(interval, symbol...)` method that validates interval and builds `kline.{interval}.{symbol}` topics, and a `KlineStreamEvent` model in `models/` mirroring the Bybit V5 push payload.

**Tech Stack:** Go 1.x, `github.com/gorilla/websocket`, existing project patterns (package `bybit_connector`, models package).

**Reference docs:** https://bybit-exchange.github.io/docs/v5/websocket/public/kline

---

### Task 1: KlineInterval type + constants

**Files:**
- Modify: `consts.go` (append at end)

- [ ] **Step 1: Add type + constants + IsValid**

```go
// KlineInterval is the websocket/REST kline interval used by Bybit V5.
// See: https://bybit-exchange.github.io/docs/v5/websocket/public/kline
type KlineInterval string

const (
	KlineInterval1m  KlineInterval = "1"
	KlineInterval3m  KlineInterval = "3"
	KlineInterval5m  KlineInterval = "5"
	KlineInterval15m KlineInterval = "15"
	KlineInterval30m KlineInterval = "30"
	KlineInterval1h  KlineInterval = "60"
	KlineInterval2h  KlineInterval = "120"
	KlineInterval4h  KlineInterval = "240"
	KlineInterval6h  KlineInterval = "360"
	KlineInterval12h KlineInterval = "720"
	KlineIntervalD   KlineInterval = "D"
	KlineIntervalW   KlineInterval = "W"
	KlineIntervalM   KlineInterval = "M"
)

func (k KlineInterval) IsValid() bool {
	switch k {
	case KlineInterval1m, KlineInterval3m, KlineInterval5m, KlineInterval15m,
		KlineInterval30m, KlineInterval1h, KlineInterval2h, KlineInterval4h,
		KlineInterval6h, KlineInterval12h, KlineIntervalD, KlineIntervalW, KlineIntervalM:
		return true
	}
	return false
}
```

---

### Task 2: KlineStreamEvent model

**Files:**
- Modify: `models/marketResponse.go` (append after existing kline structs)

- [ ] **Step 1: Add event structs**

```go
// KlineStreamItem is a single candle update pushed over the websocket
// kline.{interval}.{symbol} topic. Bybit V5 pushes an array of these
// whenever the current candle ticks or closes.
// Docs: https://bybit-exchange.github.io/docs/v5/websocket/public/kline
type KlineStreamItem struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Interval  string `json:"interval"`
	Open      string `json:"open"`
	Close     string `json:"close"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Volume    string `json:"volume"`
	Turnover  string `json:"turnover"`
	Confirm   bool   `json:"confirm"`
	Timestamp int64  `json:"timestamp"`
}

type KlineStreamEvent struct {
	Topic string            `json:"topic"`
	Type  string            `json:"type"`
	Ts    int64             `json:"ts"`
	Data  []KlineStreamItem `json:"data"`
}
```

---

### Task 3: SubscribeKline helper

**Files:**
- Modify: `bybit_websocket.go` (append after `SendSubscription`)

- [ ] **Step 1: Add helper**

```go
// SubscribeKline subscribes to the public kline stream for the given
// interval and one or more symbols. Topic format: kline.{interval}.{symbol}.
// Reference: https://bybit-exchange.github.io/docs/v5/websocket/public/kline
func (b *WebSocket) SubscribeKline(interval KlineInterval, symbols ...string) (*WebSocket, error) {
	if !interval.IsValid() {
		return b, fmt.Errorf("bybit_connector: invalid kline interval %q", string(interval))
	}
	if len(symbols) == 0 {
		return b, fmt.Errorf("bybit_connector: SubscribeKline requires at least one symbol")
	}
	topics := make([]string, 0, len(symbols))
	for _, s := range symbols {
		if s == "" {
			return b, fmt.Errorf("bybit_connector: SubscribeKline received an empty symbol")
		}
		topics = append(topics, fmt.Sprintf("kline.%s.%s", interval, s))
	}
	return b.SendSubscription(topics)
}
```

---

### Task 4: Unit tests

**Files:**
- Create: `kline_ws_test.go`

- [ ] **Step 1: Add tests for interval validation and topic building**

Test validates `KlineInterval.IsValid()` and asserts that `SubscribeKline` rejects invalid input before touching the socket (needs no real connection).

Note: `SubscribeKline` with a valid interval and a disconnected client would nil-deref `b.conn`, so tests only cover the validation branches which return early with an error.

```go
package bybit_connector

import (
	"strings"
	"testing"
)

func TestKlineInterval_IsValid(t *testing.T) {
	valid := []KlineInterval{
		KlineInterval1m, KlineInterval3m, KlineInterval5m, KlineInterval15m,
		KlineInterval30m, KlineInterval1h, KlineInterval2h, KlineInterval4h,
		KlineInterval6h, KlineInterval12h, KlineIntervalD, KlineIntervalW, KlineIntervalM,
	}
	for _, v := range valid {
		if !v.IsValid() {
			t.Errorf("expected %q to be valid", v)
		}
	}
	for _, v := range []KlineInterval{"", "2", "1d", "daily"} {
		if v.IsValid() {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}

func TestSubscribeKline_RejectsInvalidInterval(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	_, err := ws.SubscribeKline(KlineInterval("bogus"), "BTCUSDT")
	if err == nil || !strings.Contains(err.Error(), "invalid kline interval") {
		t.Fatalf("expected invalid interval error, got %v", err)
	}
}

func TestSubscribeKline_RejectsNoSymbols(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	_, err := ws.SubscribeKline(KlineInterval1m)
	if err == nil || !strings.Contains(err.Error(), "at least one symbol") {
		t.Fatalf("expected no-symbols error, got %v", err)
	}
}

func TestSubscribeKline_RejectsEmptySymbol(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	_, err := ws.SubscribeKline(KlineInterval1m, "BTCUSDT", "")
	if err == nil || !strings.Contains(err.Error(), "empty symbol") {
		t.Fatalf("expected empty-symbol error, got %v", err)
	}
}
```

- [ ] **Step 2: Run**

Run: `go test ./... -run Kline -v`
Expected: PASS

---

### Task 5: Example

**Files:**
- Create: `examples/Websocket/websocket_kline_channel.go`

- [ ] **Step 1: Write example mirroring websocket_spot_orderbook_channel.go style**

```go
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
```

---

### Task 6: Build verification

- [ ] **Step 1:** `go build ./...`
- [ ] **Step 2:** `go test ./...`
- [ ] **Step 3:** `go vet ./...`

Expected: all succeed, no new warnings.
