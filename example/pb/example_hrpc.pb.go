// Code generated by protoc-gen-hrpc. DO NOT EDIT.
package examplepb

import (
	context "context"
	hrpc "github.com/joesonw/hrpc"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
)

type ExampleClient interface {
	Echo(ctx context.Context, req *EchoMessage, opts ...hrpc.CallOption) (*EchoMessage, error)
}

func NewExampleClient(client *hrpc.Client) ExampleClient {
	return &exampleClient{
		client: client,
	}
}

type exampleClient struct {
	client *hrpc.Client
}

func (x *exampleClient) Echo(ctx context.Context, in *EchoMessage, opts ...hrpc.CallOption) (*EchoMessage, error) {
	out := &EchoMessage{}
	err := x.client.Invoke(ctx, "/example.Example/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ExampleServer interface {
	Echo(ctx context.Context, req *EchoMessage) (*EchoMessage, error)
}

func RegisterExampleServer(hrpcServer *hrpc.Server, server ExampleServer) {
	hrpcServer.RegisterService(&hrpc.ServiceDesc{
		Name:   "example.Example",
		Server: server,
		Methods: []*hrpc.MethodDesc{
			{
				Name: "Echo",
				NewRequest: func() proto.Message {
					return &EchoMessage{}
				},
				HandleFunc: func(r *http.Request, body proto.Message) (proto.Message, error) {
					return server.Echo(r.Context(), body.(*EchoMessage))
				},
			},
		},
	})
}
