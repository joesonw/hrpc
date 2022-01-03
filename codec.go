package hrpc

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Codec interface {
	Unmarshal([]byte, proto.Message) error
	Marshal(proto.Message) ([]byte, error)
}

type JSONCodec struct{}

func (JSONCodec) Unmarshal(data []byte, in proto.Message) error {
	return json.Unmarshal(data, in)
}

func (JSONCodec) Marshal(in proto.Message) ([]byte, error) {
	return json.Marshal(in)
}

type ProtoCodec struct{}

func (ProtoCodec) Unmarshal(data []byte, in proto.Message) error {
	return proto.Unmarshal(data, in)
}

func (ProtoCodec) Marshal(in proto.Message) ([]byte, error) {
	return proto.Marshal(in)
}

type ProtoJSONCodec struct {
	/* protojson.UnmarshalOptions */
	AllowPartial   bool
	DiscardUnknown bool

	/* protojson.MarshalOptions */
	Multiline       bool
	Ident           string
	UseProtoNames   bool
	UseEnumNumbers  bool
	EmitUnpopulated bool
}

func (c ProtoJSONCodec) Unmarshal(data []byte, in proto.Message) error {
	return (protojson.UnmarshalOptions{
		AllowPartial:   c.AllowPartial,
		DiscardUnknown: c.DiscardUnknown,
	}).Unmarshal(data, in)
}

func (c ProtoJSONCodec) Marshal(in proto.Message) ([]byte, error) {
	return (protojson.MarshalOptions{
		Multiline:       c.Multiline,
		Indent:          c.Ident,
		AllowPartial:    c.AllowPartial,
		UseProtoNames:   c.UseProtoNames,
		UseEnumNumbers:  c.UseEnumNumbers,
		EmitUnpopulated: c.EmitUnpopulated,
	}).Marshal(in)
}
