package protocol

import "github.com/jennal/goplay/handler/pkg"
import "bytes"
import "encoding/binary"
import "fmt"

type HeaderDecoder struct {
}

func (self HeaderDecoder) UnmarshalHeader(data []byte, header *pkg.Header) int {
	buffer := bytes.NewBuffer(data)

	b, _ := buffer.ReadByte()
	header.Type = pkg.PackageType(b)
	b, _ = buffer.ReadByte()
	header.Encoding = pkg.EncodingType(b)
	b, _ = buffer.ReadByte()
	header.ID = pkg.PackageID(b)

	fmt.Println("ContentSize", header.ContentSize, data)
	r := bytes.NewReader(data[3:7])
	binary.Read(r, binary.BigEndian, &header.ContentSize)
	fmt.Println("ContentSize", header.ContentSize, data[3:7])

	return 7
}
