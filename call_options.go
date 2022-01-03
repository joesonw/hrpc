package hrpc

type callOptions struct {
}

type CallOption func(o *callOptions)
