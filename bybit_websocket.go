package bybit_connector

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type MessageHandler func(message string) error

// WebSocket is a Bybit V5 websocket client. After construction, call
// Connect once; the client owns one read goroutine and one ping goroutine
// for its entire lifetime, transparently re-dialing on read errors.
// Call Disconnect to stop both and release the underlying connection.
//
// All exported methods are safe to call from multiple goroutines.
type WebSocket struct {
	// Immutable after construction.
	url          string
	apiKey       string
	apiSecret    string
	maxAliveTime string
	pingInterval int
	onMessage    MessageHandler
	logger       zerolog.Logger

	// Lifecycle. ctx/cancel are created in the constructor so background
	// goroutines never observe a nil context.
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// mu protects conn, isConnected, started, closed.
	mu          sync.Mutex
	conn        *websocket.Conn
	isConnected bool
	started     bool
	closed      bool

	// writeMu serialises WriteMessage calls. gorilla/websocket only
	// supports a single concurrent writer, but ping and user code may
	// both want to send.
	writeMu sync.Mutex
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

func WithLogLevel(level zerolog.Level) WebsocketOption {
	return func(c *WebSocket) {
		c.logger = c.logger.Level(level)
	}
}

func NewBybitPrivateWebSocket(url, apiKey, apiSecret string, handler MessageHandler, options ...WebsocketOption) *WebSocket {
	ctx, cancel := context.WithCancel(context.Background())
	c := &WebSocket{
		url:          url,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		pingInterval: 20,
		onMessage:    handler,
		logger:       zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		ctx:          ctx,
		cancel:       cancel,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func NewBybitPublicWebSocket(url string, handler MessageHandler) *WebSocket {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocket{
		url:          url,
		pingInterval: 20,
		onMessage:    handler,
		logger:       zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (b *WebSocket) SetMessageHandler(handler MessageHandler) {
	b.onMessage = handler
}

func (b *WebSocket) SetLogLevel(level zerolog.Level) {
	b.logger = b.logger.Level(level)
}

// Connect dials the underlying websocket and starts the read and ping
// goroutines. It is idempotent: subsequent calls on a running client are
// no-ops, calls on a closed client return an error. The first dial must
// succeed; subsequent reconnects are handled transparently by the read
// loop with exponential backoff.
func (b *WebSocket) Connect() error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return errors.New("bybit_connector: websocket is closed")
	}
	if b.started {
		b.mu.Unlock()
		return nil
	}
	b.started = true
	b.mu.Unlock()

	if err := b.dial(); err != nil {
		b.mu.Lock()
		b.started = false
		b.mu.Unlock()
		return err
	}

	b.wg.Add(2)
	go b.readLoop()
	go b.pingLoop()
	return nil
}

// dial performs a single dial+auth cycle. On failure the connection is
// closed and isConnected stays false; on success the new conn is stored
// under the lock.
func (b *WebSocket) dial() error {
	wssURL := b.url
	if b.maxAliveTime != "" {
		wssURL += "?max_alive_time=" + b.maxAliveTime
	}

	conn, _, err := websocket.DefaultDialer.Dial(wssURL, nil)
	if err != nil {
		return fmt.Errorf("bybit_connector: websocket dial: %w", err)
	}

	b.mu.Lock()
	b.conn = conn
	b.isConnected = true
	b.mu.Unlock()

	if b.requiresAuthentication() {
		if err := b.sendAuth(); err != nil {
			b.mu.Lock()
			_ = conn.Close()
			b.conn = nil
			b.isConnected = false
			b.mu.Unlock()
			return fmt.Errorf("bybit_connector: websocket auth: %w", err)
		}
	}

	b.logger.Debug().Str("url", b.url).Msg("WebSocket connected")
	return nil
}

// readLoop owns the gorilla websocket reader. On read errors it closes
// the conn, marks the client disconnected, and re-dials with exponential
// backoff capped at 30s. It exits only when the context is cancelled.
func (b *WebSocket) readLoop() {
	defer b.wg.Done()

	const (
		initialBackoff = time.Second
		maxBackoff     = 30 * time.Second
	)
	backoff := initialBackoff

	for {
		if b.ctx.Err() != nil {
			return
		}

		b.mu.Lock()
		conn := b.conn
		connected := b.isConnected
		b.mu.Unlock()

		if !connected || conn == nil {
			select {
			case <-b.ctx.Done():
				return
			case <-time.After(backoff):
			}
			if err := b.dial(); err != nil {
				b.logger.Warn().Err(err).Dur("backoff", backoff).Msg("WebSocket reconnect failed")
				if backoff < maxBackoff {
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				}
				continue
			}
			backoff = initialBackoff
			continue
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			if b.ctx.Err() != nil {
				return
			}
			b.logger.Warn().Err(err).Msg("WebSocket read error, will reconnect")
			b.mu.Lock()
			if b.conn == conn {
				_ = conn.Close()
				b.conn = nil
				b.isConnected = false
			}
			b.mu.Unlock()
			continue
		}

		if b.onMessage != nil {
			if herr := b.onMessage(string(message)); herr != nil {
				b.logger.Warn().Err(herr).Msg("WebSocket message handler returned error (continuing)")
			}
		}
	}
}

// pingLoop sends an application-level ping every pingInterval seconds.
// On send failure it nudges the read loop to reconnect by closing the
// current conn under the lock.
func (b *WebSocket) pingLoop() {
	defer b.wg.Done()

	if b.pingInterval <= 0 {
		b.logger.Debug().Int("pingInterval", b.pingInterval).Msg("Ping interval non-positive, ping loop disabled")
		return
	}
	ticker := time.NewTicker(time.Duration(b.pingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-ticker.C:
		}

		b.mu.Lock()
		connected := b.isConnected
		b.mu.Unlock()
		if !connected {
			continue
		}

		now := time.Now().Unix()
		pingMessage := map[string]string{
			"op":     "ping",
			"req_id": fmt.Sprintf("%d", now),
		}
		data, err := json.Marshal(pingMessage)
		if err != nil {
			b.logger.Warn().Err(err).Msg("Failed to marshal ping message")
			continue
		}
		if err := b.send(string(data)); err != nil {
			b.logger.Warn().Err(err).Msg("Failed to send WebSocket ping")
			b.mu.Lock()
			if b.conn != nil {
				_ = b.conn.Close()
				b.conn = nil
				b.isConnected = false
			}
			b.mu.Unlock()
			continue
		}
		b.logger.Debug().Int64("timestamp", now).Msg("WebSocket ping sent")
	}
}

// Disconnect cancels the lifetime context, closes the current
// connection, and waits for the read and ping goroutines to exit. It is
// safe to call before Connect, after a failed Connect, or multiple times.
func (b *WebSocket) Disconnect() error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil
	}
	b.closed = true
	cancel := b.cancel
	conn := b.conn
	wasStarted := b.started
	b.conn = nil
	b.isConnected = false
	b.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	var err error
	if conn != nil {
		err = conn.Close()
	}
	if wasStarted {
		b.wg.Wait()
	}
	return err
}

func (b *WebSocket) SendSubscription(args []string) (*WebSocket, error) {
	reqID := uuid.New().String()
	subMessage := map[string]interface{}{
		"req_id": reqID,
		"op":     "subscribe",
		"args":   args,
	}
	b.logger.Debug().Str("req_id", reqID).Msg("Sending WebSocket subscription")
	if err := b.sendAsJson(subMessage); err != nil {
		b.logger.Warn().Err(err).Msg("Failed to send WebSocket subscription")
		return b, err
	}
	return b, nil
}

// SubscribeKline subscribes to the public kline stream for the given
// interval and one or more symbols. The topic format is
// kline.{interval}.{symbol} (e.g. kline.1.BTCUSDT).
// Docs: https://bybit-exchange.github.io/docs/v5/websocket/public/kline
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

// SendRequest sends a custom request over the WebSocket connection.
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
	b.logger.Debug().Str("req_id", finalReqId).Str("op", op).Msg("Sending WebSocket request")
	if err := b.sendAsJson(request); err != nil {
		b.logger.Warn().Err(err).Msg("Failed to send WebSocket request")
		return b, err
	}
	return b, nil
}

func (b *WebSocket) SendTradeRequest(tradeRequest map[string]interface{}) (*WebSocket, error) {
	b.logger.Debug().Str("op", fmt.Sprintf("%v", tradeRequest["op"])).Msg("Sending WebSocket trade request")
	if err := b.sendAsJson(tradeRequest); err != nil {
		b.logger.Warn().Err(err).Msg("Failed to send WebSocket trade request")
		return b, err
	}
	return b, nil
}

func (b *WebSocket) requiresAuthentication() bool {
	return b.url == WEBSOCKET_PRIVATE_MAINNET || b.url == WEBSOCKET_PRIVATE_TESTNET ||
		b.url == WEBSOCKET_TRADE_MAINNET || b.url == WEBSOCKET_TRADE_TESTNET ||
		b.url == WEBSOCKET_TRADE_DEMO || b.url == WEBSOCKET_PRIVATE_DEMO
}

func (b *WebSocket) sendAuth() error {
	expires := time.Now().UnixNano()/1e6 + 10000
	val := fmt.Sprintf("GET/realtime%d", expires)

	h := hmac.New(sha256.New, []byte(b.apiSecret))
	h.Write([]byte(val))
	signature := hex.EncodeToString(h.Sum(nil))
	reqID := uuid.New().String()
	b.logger.Debug().Str("req_id", reqID).Int64("expires", expires).Msg("Sending WebSocket authentication")

	authMessage := map[string]interface{}{
		"req_id": reqID,
		"op":     "auth",
		"args":   []interface{}{b.apiKey, expires, signature},
	}
	return b.sendAsJson(authMessage)
}

func (b *WebSocket) sendAsJson(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.send(string(data))
}

// send serialises WriteMessage calls and returns an error rather than
// panicking when the connection is currently absent.
func (b *WebSocket) send(message string) error {
	b.mu.Lock()
	conn := b.conn
	b.mu.Unlock()
	if conn == nil {
		return errors.New("bybit_connector: websocket not connected")
	}

	b.writeMu.Lock()
	defer b.writeMu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}
