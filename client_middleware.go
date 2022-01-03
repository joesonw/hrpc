package hrpc

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

type ClientInvoker func(r *http.Request, client *http.Client, req, reply proto.Message, opts ...CallOption) (*http.Response, error)

// ClientMiddleware executed for each call on client
type ClientMiddleware func(r *http.Request, client *http.Client, method string, req, reply proto.Message, next ClientInvoker, opts ...CallOption) (*http.Response, error)

// ChainClientMiddleware chain ChainClientMiddleware array into one
func ChainClientMiddleware(middlewares ...ClientMiddleware) ClientMiddleware {
	n := len(middlewares) - 1
	return func(r *http.Request, client *http.Client, method string, req, reply proto.Message, next ClientInvoker, opts ...CallOption) (*http.Response, error) {
		_next := next
		for i := n; i >= 0; i-- {
			m := middlewares[i]
			n := _next
			_next = func(r *http.Request, client *http.Client, req, reply proto.Message, opts ...CallOption) (*http.Response, error) {
				return m(r, client, method, req, reply, n, opts...)
			}
		}

		return _next(r, client, req, reply, opts...)
	}
}
