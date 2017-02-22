package protocol

import "github.com/jennal/goplay/handler/pkg"

type Encoder interface {
	Marshal(obj *pkg.Header, content interface{}) ([]byte, error)
	MarshalHeader(header *pkg.Header) ([]byte, error)
	MarshalContent(content interface{}) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, header *pkg.Header, content interface{}) error
	UnmarshalHeader(data []byte, header *pkg.Header) (int, error)
	UnmarshalContent(data []byte, content interface{}) error
}

type EncodeDecoder interface {
	Encoder
	Decoder
}
