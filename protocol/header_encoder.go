package protocol

import "github.com/jennal/goplay/handler/pkg"
import "encoding/binary"
import "bytes"

type HeaderEncoder struct {
}

func (self HeaderEncoder) MarshalHeader(header *pkg.Header) []byte {
	var buffer bytes.Buffer

	buffer.Write([]byte{
		byte(header.Type),
		byte(header.Encoding),
		byte(header.ID),
	})
	binary.Write(&buffer, binary.BigEndian, header.ContentSize)

	return buffer.Bytes()
}
