package protocol

import "github.com/jennal/goplay/handler/pkg"

type Encoder interface {
	Marshal(obj *pkg.Package) []byte
	MarshalContent(obj interface{}) []byte
}

type Decoder interface {
	Unmarshal(data []byte) *pkg.Package
	UnmarshalContent(data []byte) interface{}
}

type EncodeDecoder interface {
	Encoder
	Decoder
}
