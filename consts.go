package bybit_connector

const (
	Name    = "bybit.api.go"
	Version = "0.1.0"
	// Https
	MAINNET        = "https://api.bybit.com"
	MAINNET_BACKT  = "https://api.bytick.com"
	TESTNET        = "https://api-testnet.bybit.com"
	DEMO_ENV       = "https://api-demo.bybit.com"
	NETHERLAND_ENV = "https://api.bybit.nl"
	HONGKONG_ENV   = "https://api.byhkbit.com"
	TURKEY_ENV     = "https://api.bybit-tr.com"
	KAZAKHSTAN_ENV = "https://api.bybit.kz"

	// WebSocket public channel - Mainnet
	SPOT_MAINNET    = "wss://stream.bybit.com/v5/public/spot"
	LINEAR_MAINNET  = "wss://stream.bybit.com/v5/public/linear"
	INVERSE_MAINNET = "wss://stream.bybit.com/v5/public/inverse"
	SPREAD_MAINNET  = "wss://stream.bybit.com/v5/public/spread"

	// WebSocket public channel - Testnet
	SPOT_TESTNET    = "wss://stream-testnet.bybit.com/v5/public/spot"
	LINEAR_TESTNET  = "wss://stream-testnet.bybit.com/v5/public/linear"
	INVERSE_TESTNET = "wss://stream-testnet.bybit.com/v5/public/inverse"
	OPTION_TESTNET  = "wss://stream-testnet.bybit.com/v5/public/option"
	SPREAD_TESTNET  = "wss://stream-testnet.bybit.com/v5/public/spread"

	// WebSocket private channel
	WEBSOCKET_PRIVATE_MAINNET = "wss://stream.bybit.com/v5/private"
	WEBSOCKET_TRADE_MAINNET   = "wss://stream.bybit.com/v5/trade"
	WEBSOCKET_PRIVATE_TESTNET = "wss://stream-testnet.bybit.com/v5/private"
	WEBSOCKET_TRADE_TESTNET   = "wss://stream-testnet.bybit.com/v5/trade"
	WEBSOCKET_PRIVATE_DEMO    = "wss://stream-demo.bybit.com/v5/private"
	WEBSOCKET_TRADE_DEMO      = "wss://stream-demo.bybit.com/v5/trade"

	// Globals
	TimestampKey  = "X-BAPI-TIMESTAMP"
	SignatureKey  = "X-BAPI-SIGN"
	ApiRequestKey = "X-BAPI-API-KEY"
	RecvWindowKey = "X-BAPI-RECV-WINDOW"
	SignTypeKey   = "X-BAPI-SIGN-TYPE"
)

// KlineInterval is a Bybit V5 kline interval, used for both REST and the
// public websocket kline.{interval}.{symbol} stream.
// Docs: https://bybit-exchange.github.io/docs/v5/websocket/public/kline
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
