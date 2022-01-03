all:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install github.com/joesonw/hrpc/cmd/protoc-gen-hrpc
	protoc \
		--proto_path=./ \
		--hrpc_out=./ \
		--hrpc_opt=module=example \
		--go_out=./ \
		--go_opt=module=example \
	 	./example.proto