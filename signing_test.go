package bybit_connector

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// TestComputeSignature_GoldenVectors locks the Bybit V5 REST signing
// algorithm: HMAC-SHA256 over the literal concatenation of
// (timestamp || apiKey || recvWindow || payload), hex-encoded, with the
// API secret as the HMAC key. Any change to the order, casing, or
// hashing primitive will break these vectors and surface immediately.
func TestComputeSignature_GoldenVectors(t *testing.T) {
	const (
		apiSecret  = "mySuperSecret123"
		apiKey     = "myKey"
		timestamp  = "1700000000000"
		recvWindow = "5000"
	)

	cases := []struct {
		name    string
		payload string
		want    string
	}{
		{
			name:    "POST body",
			payload: `{"category":"linear","symbol":"BTCUSDT","side":"Buy","orderType":"Market","qty":"0.001"}`,
			want:    "da8fcd450c9ba4ef10e0ab5ad6c4670d51d4009a76b3b4bbfe55325434b97151",
		},
		{
			name:    "GET query string",
			payload: "category=linear&symbol=BTCUSDT",
			want:    "25c1a58244fed638faa9a590d2a023b3f7e67d22cac24fbc24cab1f6eb1a6494",
		},
		{
			name:    "empty payload",
			payload: "",
			want:    "c4529dd020527029152920f689df7cb61d6b6c0f4ff7104f127c76daa50903e7",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := computeSignature(apiSecret, apiKey, timestamp, recvWindow, tc.payload)
			if got != tc.want {
				t.Errorf("computeSignature() = %s, want %s", got, tc.want)
			}
		})
	}
}

// TestComputeSignature_MatchesStdlib cross-checks the helper against a
// raw stdlib reimplementation. This catches subtle mutations to the
// helper such as silent payload truncation or a forgotten field.
func TestComputeSignature_MatchesStdlib(t *testing.T) {
	const (
		apiSecret  = "anotherSecret"
		apiKey     = "anotherKey"
		timestamp  = "1745000000000"
		recvWindow = "8000"
		payload    = `{"foo":"bar","baz":42}`
	)
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(timestamp + apiKey + recvWindow + payload))
	want := hex.EncodeToString(h.Sum(nil))

	got := computeSignature(apiSecret, apiKey, timestamp, recvWindow, payload)
	if got != want {
		t.Errorf("computeSignature() = %s, want %s", got, want)
	}
}

// TestParseRequest_SignsPostRequest exercises the full parseRequest
// path: it builds a signed POST Request, runs parseRequest through a
// real Client, and verifies that the resulting X-BAPI-SIGN header
// matches what computeSignature would have produced for the same
// timestamp the function chose. This locks the wiring between
// parseRequest and computeSignature.
func TestParseRequest_SignsPostRequest(t *testing.T) {
	c := NewBybitHttpClient("apiKey", "apiSecret", WithBaseURL("https://example.invalid"))

	body := []byte(`{"category":"linear","symbol":"BTCUSDT"}`)
	r := &Request{
		method:   "POST",
		endpoint: "/v5/order/create",
		secType:  secTypeSigned,
		Params:   body,
	}
	if err := c.parseRequest(r); err != nil {
		t.Fatalf("parseRequest: %v", err)
	}

	timestamp := r.header.Get(TimestampKey)
	recvWindow := r.header.Get(RecvWindowKey)
	apiKeyHdr := r.header.Get(ApiRequestKey)
	gotSig := r.header.Get(SignatureKey)

	if apiKeyHdr != "apiKey" {
		t.Errorf("X-BAPI-API-KEY = %q, want %q", apiKeyHdr, "apiKey")
	}
	if recvWindow != "5000" {
		t.Errorf("X-BAPI-RECV-WINDOW = %q, want %q", recvWindow, "5000")
	}
	if timestamp == "" {
		t.Fatal("X-BAPI-TIMESTAMP is empty")
	}

	want := computeSignature("apiSecret", "apiKey", timestamp, recvWindow, string(body))
	if gotSig != want {
		t.Errorf("X-BAPI-SIGN = %s, want %s", gotSig, want)
	}
}
