package protocol

import "github.com/jennal/goplay/handler/pkg"
import "bytes"

type BaseEncoder struct {
}

func (self BaseEncoder) Marshal(obj *pkg.Package) []byte {
	var buffer bytes.Buffer

	buffer.WriteByte(byte(obj.Type))
	buffer.WriteByte(byte(obj.Encoding))
	buffer.WriteByte(byte(obj.ID))

	contentBuff := self.MarshalContent(obj.Content)
	buffer.Write(contentBuff)

	return buffer.Bytes()
}

func (self BaseEncoder) MarshalContent(obj interface{}) []byte {
	return []byte{}
}
