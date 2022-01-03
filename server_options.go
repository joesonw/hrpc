package hrpc

type serverOptions struct {
	codec      Codec
	middleware ServerMiddleware
}

type ServerOption func(o *serverOptions)

func WithServerCodec(codec Codec) ServerOption {
	return func(o *serverOptions) {
		o.codec = codec
	}
}

func WithServerMiddleware(middleware ServerMiddleware) ServerOption {
	return func(o *serverOptions) {
		o.middleware = middleware
	}
}
