package hrpc

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/joesonw/hrpc/status"
)

type ServiceDesc struct {
	Name    string
	Server  interface{}
	Methods []*MethodDesc
}

type MethodDesc struct {
	Name       string
	NewRequest func() proto.Message
	HandleFunc func(r *http.Request, body proto.Message) (proto.Message, error)
}

type methodDesc struct {
	*MethodDesc
	ServiceDesc *ServiceDesc
}

type Server struct {
	services map[string]*ServiceDesc
	methods  map[string]*methodDesc

	codec      Codec
	middleware ServerMiddleware
}

func NewServer(options ...ServerOption) *Server {
	o := &serverOptions{
		codec: JSONCodec{},
	}

	for _, apply := range options {
		apply(o)
	}

	return &Server{
		services: map[string]*ServiceDesc{},
		methods:  map[string]*methodDesc{},

		codec:      o.codec,
		middleware: o.middleware,
	}
}

func (s *Server) RegisterService(desc *ServiceDesc) {
	s.services[desc.Name] = desc
	for _, method := range desc.Methods {
		s.methods["/"+desc.Name+"/"+method.Name] = &methodDesc{
			MethodDesc:  method,
			ServiceDesc: desc,
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost || r.Body == nil {
		return
	}

	fullMethod := r.URL.Path
	method, ok := s.methods[fullMethod]
	if !ok {
		http.Error(w, fmt.Sprintf("Method %q not found", fullMethod), http.StatusNotFound)
		return
	}

	b, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := method.NewRequest()
	if err := s.codec.Unmarshal(b, req); err != nil {
		http.Error(w, err.Error(), status.HTTPStatus(err))
		return
	}

	var res proto.Message
	if s.middleware != nil {
		res, err = s.middleware(r, req, &ServerMiddlewareInfo{
			Server:     method.ServiceDesc.Server,
			FullMethod: fullMethod,
		}, method.HandleFunc)
	} else {
		res, err = method.HandleFunc(r, req)
	}

	if err != nil {
		http.Error(w, err.Error(), status.HTTPStatus(err))
		return
	}

	b, err = s.codec.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), status.HTTPStatus(err))
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(b)
}
