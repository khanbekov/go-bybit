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

func TestWebSocket_DisconnectBeforeConnect_NoPanic(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	if err := ws.Disconnect(); err != nil {
		t.Fatalf("Disconnect on fresh client should not error, got %v", err)
	}
}

func TestWebSocket_DisconnectIsIdempotent(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	if err := ws.Disconnect(); err != nil {
		t.Fatalf("first Disconnect: %v", err)
	}
	if err := ws.Disconnect(); err != nil {
		t.Fatalf("second Disconnect: %v", err)
	}
}

func TestWebSocket_ConnectAfterDisconnect_Errors(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	_ = ws.Disconnect()
	if err := ws.Connect(); err == nil || !strings.Contains(err.Error(), "closed") {
		t.Fatalf("expected closed error, got %v", err)
	}
}

func TestWebSocket_ConnectFailsCleanly(t *testing.T) {
	ws := NewBybitPublicWebSocket("ws://127.0.0.1:1", nil)
	err := ws.Connect()
	if err == nil {
		t.Fatal("expected dial error against unreachable address")
	}
	if !strings.Contains(err.Error(), "websocket dial") {
		t.Fatalf("expected wrapped dial error, got %v", err)
	}
	// After failed Connect, Disconnect must still be safe.
	if derr := ws.Disconnect(); derr != nil {
		t.Fatalf("Disconnect after failed Connect: %v", derr)
	}
}

func TestSendBeforeConnect_Errors(t *testing.T) {
	ws := NewBybitPublicWebSocket("wss://example.invalid", nil)
	_, err := ws.SendSubscription([]string{"orderbook.1.BTCUSDT"})
	if err == nil || !strings.Contains(err.Error(), "not connected") {
		t.Fatalf("expected not-connected error, got %v", err)
	}
}
