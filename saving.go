package bybit_connector

import (
	"context"
	"github.com/khanbekov/go-bybit/handlers"
	"net/http"
)

func (s *BybitClientRequest) GetEarnProductInfo(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	if err = handlers.ValidateParams(s.Params); err != nil {
		return nil, err
	}
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/earn/product",
		secType:  secTypeSigned,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) RedeemEarnOrder(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodPost,
		endpoint: "/v5/earn/place-order",
		secType:  secTypeSigned,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetEarnRedeemOrder(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/earn/order",
		secType:  secTypeSigned,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}

func (s *BybitClientRequest) GetEarnRedeemPosition(ctx context.Context, opts ...RequestOption) (res *ServerResponse, err error) {
	r := &Request{
		method:   http.MethodGet,
		endpoint: "/v5/earn/position",
		secType:  secTypeSigned,
	}
	data, err := SendRequest(ctx, opts, r, s, err)
	return GetServerResponse(err, data)
}
