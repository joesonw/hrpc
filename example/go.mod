module example

go 1.17

require (
	github.com/joesonw/hrpc v0.0.1
	google.golang.org/protobuf v1.27.1
)

require github.com/joesonw/proto-tools v0.1.3 // indirect

replace github.com/joesonw/hrpc => ../
