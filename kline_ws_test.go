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
