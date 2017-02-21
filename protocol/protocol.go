package protocol

import "github.com/jennal/goplay/handler/pkg"

type Encoder interface {
	Marshal(obj *pkg.Header, content interface{}) []byte
	MarshalHeader(header *pkg.Header) []byte
	MarshalContent(content interface{}) []byte
}

type Decoder interface {
	Unmarshal(data []byte, header *pkg.Header, content interface{})
	UnmarshalHeader(data []byte, header *pkg.Header) int
	UnmarshalContent(data []byte, content interface{})
}

type EncodeDecoder interface {
	Encoder
	Decoder
}
