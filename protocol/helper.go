package protocol

import (
	"bytes"

	"github.com/jennal/goplay/handler/pkg"
)

func Marshal(self Encoder, header *pkg.Header, content interface{}) ([]byte, error) {
	contentBuff, err := self.MarshalContent(content)
	if err != nil {
		return nil, err
	}

	header.ContentSize = uint32(len(contentBuff))
	headerBuff, err := self.MarshalHeader(header)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.Write(headerBuff)
	buffer.Write(contentBuff)

	return buffer.Bytes(), nil
}

func Unmarshal(self Decoder, data []byte, header *pkg.Header, content interface{}) error {
	n, err := self.UnmarshalHeader(data, header)
	if err != nil {
		return err
	}
	return self.UnmarshalContent(data[n:], content)
}
