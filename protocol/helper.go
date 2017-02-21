package protocol

import (
	"bytes"

	"github.com/jennal/goplay/handler/pkg"
)

func Marshal(self Encoder, header *pkg.Header, content interface{}) []byte {
	contentBuff := self.MarshalContent(content)
	header.ContentSize = uint32(len(contentBuff))
	headerBuff := self.MarshalHeader(header)

	var buffer bytes.Buffer
	buffer.Write(headerBuff)
	buffer.Write(contentBuff)

	return buffer.Bytes()
}

func Unmarshal(self Decoder, data []byte, header *pkg.Header, content interface{}) {
	n := self.UnmarshalHeader(data, header)
	self.UnmarshalContent(data[n:], content)
}
