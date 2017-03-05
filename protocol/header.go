package protocol

import (
	"bytes"

	"errors"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/helpers"
)

const HEADER_STATIC_SIZE = 5

type HeaderEncoder struct {
}

func (self HeaderEncoder) MarshalHeader(header *pkg.Header) ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write([]byte{
		byte(header.Type),
		byte(header.Encoding),
		byte(header.ID),
	})
	buf, err := helpers.GetBytes(header.ContentSize)
	buffer.Write(buf)
	buffer.WriteByte(byte(len(header.Route)))
	buffer.Write([]byte(header.Route))

	return buffer.Bytes(), err
}

type HeaderDecoder struct {
}

func (self HeaderDecoder) UnmarshalHeader(data []byte, header *pkg.Header) (int, error) {
	if len(data) < HEADER_STATIC_SIZE {
		return -1, errors.New("data size < HEADER_STATIC_SIZE")
	}

	buffer := bytes.NewBuffer(data)

	b, err := buffer.ReadByte()
	header.Type = pkg.PackageType(b)
	b, err = buffer.ReadByte()
	header.Encoding = pkg.EncodingType(b)
	b, err = buffer.ReadByte()
	header.ID = pkg.PackageIDType(b)

	size, err := helpers.ToUInt16(data[3:HEADER_STATIC_SIZE])
	header.ContentSize = pkg.PackageSizeType(size)
	for i := 3; i < HEADER_STATIC_SIZE; i++ {
		buffer.ReadByte()
	}
	b, err = buffer.ReadByte()
	routeSize := int(b)
	route := make([]byte, b)
	for i := 0; i < routeSize; i++ {
		route[i], err = buffer.ReadByte()
	}
	header.Route = string(route)

	return HEADER_STATIC_SIZE + 1 + routeSize, err
}
