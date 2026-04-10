bybit.go.api 0.1.0

BREAKING
WebSocket Connect() now returns error instead of *WebSocket.
Migration:
    // before
    _, _ = ws.Connect().SendSubscription([]string{"order"})
    // after
    if err := ws.Connect(); err != nil { /* handle */ }
    _, _ = ws.SendSubscription([]string{"order"})

Websocket
Add SubscribeKline helper for the public kline.{interval}.{symbol} stream
Add KlineInterval type with the 13 valid Bybit V5 intervals
Add KlineStreamItem and KlineStreamEvent models for parsing push frames
Reworked WebSocket lifecycle: race-free, idempotent Disconnect, exponential
backoff reconnect inside a single supervisor, no more goroutine leaks on
reconnect, no more nil-conn or nil-ctx panics

Bug fixes
GetIndexPriceKline pointed at /v5/market/mark-price-kline (now correct)
GetPremiumIndexPriceKline pointed at /v5/market/mark-price-kline (now correct)
SetSpotHedgeMode used GET, must be POST
RequestTestFund used GET, must be POST
LotSizeFilter.MinOrderAmt struct tag was jsoN: (now json:)
setParams no longer calls log.Fatal on marshal failure; the error is
propagated through SendRequest

Security
sendAuth no longer logs the WebSocket HMAC signature or the auth args
parseRequest / callAPI no longer dump the full *http.Request /
*http.Response (which include X-BAPI-API-KEY and X-BAPI-SIGN headers)
at debug level

Tests
Extracted computeSignature helper and locked the REST signing algorithm
with three golden vectors plus a stdlib cross-check
Added 6 lifecycle tests for the WebSocket client (idempotent Disconnect,
Connect-failure cleanup, send before Connect, Connect after Disconnect)
All tests pass under go test -race

Cleanup
Removed 4 unused free-function kline parsers (one with a broken
double-pointer return signature)
Removed deprecated V3_* websocket constants
Removed unreachable duplicate-key check in handlers.ValidateParams
Renamed SendTradeRequest parameter tradeTruest -> tradeRequest

Previous releases
---

bybit.go.api 1.0.7

Rest
Get closed option positions
Get price limit
Pre check order

Websocket
Order price limit
Insurance pool
Order price limit
Rpi order book

New crypto loan
Get borrowable coins
Get collateral coins
Get max allowed collateral reduction amount
Adjust collateral amount
Get collateral adjustment history
Get crypto loan position

Flexible loan
Borrow
Repay
Get flexible loan
Get borrow history
Get repay history

Fixed loan
Get supplying market
Get borrowing market
Create borrow order
Create supply order
Cancel borrow order
Cancel supply order
Get borrow contact info
Get supply contact info
Get borrow order info
Get supply order info
Repay
Get repayment history

Deprecated leverage token endpoints and legacy crypto loan endpoints 