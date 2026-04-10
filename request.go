package bybit_connector

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type secType int

const (
	secTypeNone   secType = iota
	secTypeSigned         // private request
)

type Params map[string]interface{}

// request define an API request
type Request struct {
	method     string
	endpoint   string
	query      url.Values
	recvWindow string
	secType    secType
	header     http.Header
	Params     []byte
	fullURL    string
	body       io.Reader
}

// addParam add param with key/value to query string
func (r *Request) addParam(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Add(key, fmt.Sprintf("%v", value))
	return r
}

// setParam set param with key/value to query string
func (r *Request) setParam(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// setParams sets the given params as either query string (GET) or JSON body (POST).
// It returns an error if the POST body cannot be marshaled — callers must
// never proceed to transport when this fails, otherwise the request would
// be sent without its payload.
func (r *Request) setParams(m Params) error {
	switch r.method {
	case http.MethodGet:
		for k, v := range m {
			r.setParam(k, v)
		}
	case http.MethodPost:
		jsonData, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("bybit_connector: marshal request params: %w", err)
		}
		r.Params = jsonData
	}
	return nil
}

func (r *Request) validate() (err error) {
	if r.query == nil {
		r.query = url.Values{}
	}
	return nil
}

// WithRecvWindow Append `WithRecvWindow(insert_recvWindow)` to request to modify the default recvWindow value
func WithRecvWindow(recvWindow string) RequestOption {
	return func(r *Request) {
		r.recvWindow = recvWindow
	}
}

// RequestOption define option type for request
type RequestOption func(*Request)
