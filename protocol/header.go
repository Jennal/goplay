package protocol

import (
	"bytes"

	"errors"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/helpers"
)

const HEADER_SIZE = 7

type HeaderEncoder struct {
}

func (self HeaderEncoder) MarshalHeader(header *pkg.Header) ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write([]byte{
		byte(header.Type),
		byte(header.Encoding),
		byte(header.ID),
	})
	buf, err := helpers.UInt32(header.ContentSize).GetBytes()
	buffer.Write(buf)

	return buffer.Bytes(), err
}

type HeaderDecoder struct {
}

func (self HeaderDecoder) UnmarshalHeader(data []byte, header *pkg.Header) (int, error) {
	if len(data) < HEADER_SIZE {
		return -1, errors.New("data size < HEADER_SIZE")
	}

	buffer := bytes.NewBuffer(data)

	b, _ := buffer.ReadByte()
	header.Type = pkg.PackageType(b)
	b, _ = buffer.ReadByte()
	header.Encoding = pkg.EncodingType(b)
	b, _ = buffer.ReadByte()
	header.ID = pkg.PackageID(b)

	var err error
	header.ContentSize, err = helpers.Bytes(data[3:7]).ToInt()

	return 7, err
}
