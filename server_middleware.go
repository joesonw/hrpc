package hrpc

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

type ServerHandler func(r *http.Request, req proto.Message) (proto.Message, error)

type ServerMiddlewareInfo struct {
	Server     interface{}
	FullMethod string
}

type ServerMiddleware func(r *http.Request, req proto.Message, info *ServerMiddlewareInfo, next ServerHandler) (proto.Message, error)

func ChainServerMiddleware(middlewares ...ServerMiddleware) ServerMiddleware {
	n := len(middlewares) - 1
	return func(r *http.Request, req proto.Message, info *ServerMiddlewareInfo, next ServerHandler) (proto.Message, error) {
		_next := next
		for i := n; i >= 0; i-- {
			m := middlewares[i]
			n := _next
			_next = func(r *http.Request, req proto.Message) (proto.Message, error) {
				return m(r, req, info, n)
			}
		}
		return _next(r, req)
	}
}
