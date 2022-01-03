package hrpc

import (
	"net/http"
)

type clientOptions struct {
	httpClient *http.Client
	codec      Codec
	middleware ClientMiddleware
}

// ClientOption apply client option
type ClientOption func(o *clientOptions)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

func WithClientCodec(codec Codec) ClientOption {
	return func(o *clientOptions) {
		o.codec = codec
	}
}

func WithClientMiddleware(middleware ClientMiddleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = middleware
	}
}
