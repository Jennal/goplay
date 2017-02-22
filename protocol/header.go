package protocol

import (
	"bytes"
	"encoding/binary"

	"errors"

	"github.com/jennal/goplay/handler/pkg"
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
	err := binary.Write(&buffer, binary.BigEndian, header.ContentSize)

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

	// fmt.Println("ContentSize", header.ContentSize, data)
	r := bytes.NewReader(data[3:7])
	err := binary.Read(r, binary.BigEndian, &header.ContentSize)
	// fmt.Println("ContentSize", header.ContentSize, data[3:7])

	return 7, err
}
