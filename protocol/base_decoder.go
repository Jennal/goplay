package protocol

import "github.com/jennal/goplay/handler/pkg"
import "bytes"

type BaseDecoder struct {
}

func (self BaseDecoder) Unmarshal(data []byte) *pkg.Package {
	result := &pkg.Package{}

	buffer := bytes.NewBuffer(data)
	b, _ := buffer.ReadByte()
	result.Type = pkg.PackageType(b)
	b, _ = buffer.ReadByte()
	result.Encoding = pkg.EncodingType(b)
	b, _ = buffer.ReadByte()
	result.ID = pkg.PackageID(b)

	result.Content = self.UnmarshalContent(data[3:])

	return result
}

func (self BaseDecoder) UnmarshalContent(data []byte) interface{} {
	return nil
}
