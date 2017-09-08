package interfaces

import "github.com/jennal/goplay/pkg"

type PackageEncoder interface {
	MarshalData(header *pkg.Header, content interface{}) ([]byte, error)
}

type PackageDecoder interface {
	UnmarshalData(header *pkg.Header, data []byte, content interface{}) error
}

type PackageEncodeDecoder interface {
	PackageEncoder
	PackageDecoder
}
