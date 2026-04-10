package bybit_connector

import (
	"context"
	"net/http"
)

func (s *BybitClientRequest) GetServerTime(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/time",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetMarketKline(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/kline",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetMarkPriceKline(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/mark-price-kline",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetIndexPriceKline(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/index-price-kline",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetPremiumIndexPriceKline(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/premium-index-price-kline",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetInstrumentInfo(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/instruments-info",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetOrderBookInfo(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/orderbook",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetMarketTickers(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/tickers",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetFundingRateHistory(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/funding/history",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetPublicRecentTrades(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/recent-trade",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetOpenInterests(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/open-interest",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetHistoryVolatility(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/historical-volatility",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetMarketInsurance(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/insurance",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetMarketRiskLimits(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/risk-limit",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetDeliveryPrice(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/delivery-price",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetLongShortRatio(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/account-ratio",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetOrderPriceLimit(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/market/price-limit",
		secType:  secTypeNone,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}
