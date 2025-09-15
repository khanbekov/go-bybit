package bybit_connector

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type MessageHandler func(message string) error

func (b *WebSocket) handleIncomingMessages() {
	for {
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			b.logger.Debug().Err(err).Msg("Error reading WebSocket message")
			b.isConnected = false
			return
		}

		if b.onMessage != nil {
			err := b.onMessage(string(message))
			if err != nil {
				b.logger.Debug().Err(err).Str("message", string(message)).Msg("Error handling WebSocket message")
				return
			}
		}
	}
}

func (b *WebSocket) monitorConnection() {
	ticker := time.NewTicker(time.Second * 5) // Check every 5 seconds
	defer ticker.Stop()

	for {
		<-ticker.C
		if !b.isConnected && b.ctx.Err() == nil { // Check if disconnected and context not done
			b.logger.Debug().Msg("Attempting to reconnect WebSocket")
			con := b.Connect() // Example, adjust parameters as needed
			if con == nil {
				b.logger.Debug().Msg("WebSocket reconnection failed")
			} else {
				b.isConnected = true
				go b.handleIncomingMessages() // Restart message handling
			}
		}

		select {
		case <-b.ctx.Done():
			return // Stop the routine if context is done
		default:
		}
	}
}

func (b *WebSocket) SetMessageHandler(handler MessageHandler) {
	b.onMessage = handler
}

type WebSocket struct {
	conn         *websocket.Conn
	url          string
	apiKey       string
	apiSecret    string
	maxAliveTime string
	pingInterval int
	onMessage    MessageHandler
	ctx          context.Context
	cancel       context.CancelFunc
	isConnected  bool
	logger       zerolog.Logger
}

type WebsocketOption func(*WebSocket)

func WithPingInterval(pingInterval int) WebsocketOption {
	return func(c *WebSocket) {
		c.pingInterval = pingInterval
	}
}

func WithMaxAliveTime(maxAliveTime string) WebsocketOption {
	return func(c *WebSocket) {
		c.maxAliveTime = maxAliveTime
	}
}

func NewBybitPrivateWebSocket(url, apiKey, apiSecret string, handler MessageHandler, options ...WebsocketOption) *WebSocket {
	c := &WebSocket{
		url:          url,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		maxAliveTime: "",
		pingInterval: 20,
		onMessage:    handler,
		logger:       zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
	}

	// Apply the provided options
	for _, opt := range options {
		opt(c)
	}

	return c
}

func NewBybitPublicWebSocket(url string, handler MessageHandler) *WebSocket {
	c := &WebSocket{
		url:          url,
		pingInterval: 20, // default is 20 seconds
		onMessage:    handler,
		logger:       zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
	}

	return c
}

func (b *WebSocket) Connect() *WebSocket {
	var err error
	wssUrl := b.url
	if b.maxAliveTime != "" {
		wssUrl += "?max_alive_time=" + b.maxAliveTime
	}
	b.conn, _, err = websocket.DefaultDialer.Dial(wssUrl, nil)

	if b.requiresAuthentication() {
		if err = b.sendAuth(); err != nil {
			b.logger.Debug().Err(err).Msg("Failed WebSocket authentication")
			return nil
		}
	}
	b.isConnected = true

	go b.handleIncomingMessages()
	go b.monitorConnection()

	b.ctx, b.cancel = context.WithCancel(context.Background())
	go ping(b)

	return b
}

func (b *WebSocket) SendSubscription(args []string) (*WebSocket, error) {
	reqID := uuid.New().String()
	subMessage := map[string]interface{}{
		"req_id": reqID,
		"op":     "subscribe",
		"args":   args,
	}
	b.logger.Debug().Interface("args", subMessage["args"]).Msg("Sending WebSocket subscription")
	if err := b.sendAsJson(subMessage); err != nil {
		b.logger.Debug().Err(err).Msg("Failed to send WebSocket subscription")
		return b, err
	}
	b.logger.Debug().Msg("WebSocket subscription sent successfully")
	return b, nil
}

// SendRequest sendRequest sends a custom request over the WebSocket connection.
func (b *WebSocket) SendRequest(op string, args map[string]interface{}, headers map[string]string, reqId ...string) (*WebSocket, error) {
	finalReqId := uuid.New().String()
	if len(reqId) > 0 && reqId[0] != "" {
		finalReqId = reqId[0]
	}

	request := map[string]interface{}{
		"reqId":  finalReqId,
		"header": headers,
		"op":     op,
		"args":   []interface{}{args},
	}
	b.logger.Debug().Interface("headers", request["header"]).Str("op", fmt.Sprintf("%v", request["op"])).Interface("args", request["args"]).Msg("Sending WebSocket request")
	if err := b.sendAsJson(request); err != nil {
		b.logger.Debug().Err(err).Msg("Failed to send WebSocket trade request")
		return b, err
	}
	b.logger.Debug().Msg("Successfully sent WebSocket trade request")
	return b, nil
}

func (b *WebSocket) SendTradeRequest(tradeTruest map[string]interface{}) (*WebSocket, error) {
	b.logger.Debug().Interface("headers", tradeTruest["header"]).Str("op", fmt.Sprintf("%v", tradeTruest["op"])).Interface("args", tradeTruest["args"]).Msg("Sending WebSocket trade request")
	if err := b.sendAsJson(tradeTruest); err != nil {
		b.logger.Debug().Err(err).Msg("Failed to send WebSocket trade request")
		return b, err
	}
	b.logger.Debug().Msg("Successfully sent WebSocket trade request")
	return b, nil
}

func ping(b *WebSocket) {
	if b.pingInterval <= 0 {
		b.logger.Debug().Int("pingInterval", b.pingInterval).Msg("Ping interval is set to a non-positive value")
		return
	}

	ticker := time.NewTicker(time.Duration(b.pingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentTime := time.Now().Unix()
			pingMessage := map[string]string{
				"op":     "ping",
				"req_id": fmt.Sprintf("%d", currentTime),
			}
			jsonPingMessage, err := json.Marshal(pingMessage)
			if err != nil {
				b.logger.Debug().Err(err).Msg("Failed to marshal ping message")
				continue
			}
			if err := b.conn.WriteMessage(websocket.TextMessage, jsonPingMessage); err != nil {
				b.logger.Debug().Err(err).Msg("Failed to send WebSocket ping")
				return
			}
			b.logger.Debug().Int64("timestamp", currentTime).Msg("WebSocket ping sent")

		case <-b.ctx.Done():
			b.logger.Debug().Msg("WebSocket ping context closed, stopping ping")
			return
		}
	}
}

func (b *WebSocket) Disconnect() error {
	b.cancel()
	b.isConnected = false
	return b.conn.Close()
}

func (b *WebSocket) requiresAuthentication() bool {
	return b.url == WEBSOCKET_PRIVATE_MAINNET || b.url == WEBSOCKET_PRIVATE_TESTNET ||
		b.url == WEBSOCKET_TRADE_MAINNET || b.url == WEBSOCKET_TRADE_TESTNET ||
		b.url == WEBSOCKET_TRADE_DEMO || b.url == WEBSOCKET_PRIVATE_DEMO
	// v3 offline
	/*
		b.url == V3_CONTRACT_PRIVATE ||
			b.url == V3_UNIFIED_PRIVATE ||
			b.url == V3_SPOT_PRIVATE
	*/
}

func (b *WebSocket) sendAuth() error {
	// Get current Unix time in milliseconds
	expires := time.Now().UnixNano()/1e6 + 10000
	val := fmt.Sprintf("GET/realtime%d", expires)

	h := hmac.New(sha256.New, []byte(b.apiSecret))
	h.Write([]byte(val))

	// Convert to hexadecimal instead of base64
	signature := hex.EncodeToString(h.Sum(nil))
	b.logger.Debug().Str("signature", signature).Msg("Generated WebSocket authentication signature")

	authMessage := map[string]interface{}{
		"req_id": uuid.New(),
		"op":     "auth",
		"args":   []interface{}{b.apiKey, expires, signature},
	}
	b.logger.Debug().Interface("args", authMessage["args"]).Msg("Sending WebSocket authentication")
	return b.sendAsJson(authMessage)
}

func (b *WebSocket) sendAsJson(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.send(string(data))
}

func (b *WebSocket) send(message string) error {
	return b.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
