package protocol

import (
	"bytes"
	"encoding/gob"

	"github.com/jennal/goplay/handler/pkg"
)

/* TODO: */
type ProtobufDecoder struct {
	HeaderDecoder
}

func (self ProtobufDecoder) Unmarshal(data []byte, header *pkg.Header, content interface{}) {
	Unmarshal(self, data, header, content)
}

func (self ProtobufDecoder) UnmarshalContent(data []byte, content interface{}) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(content)
}
