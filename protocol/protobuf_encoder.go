package protocol

import (
	"bytes"
	"encoding/gob"

	"github.com/jennal/goplay/handler/pkg"
)

/* TODO: */
type ProtobufEncoder struct {
	HeaderEncoder
}

func (self ProtobufEncoder) Marshal(header *pkg.Header, content interface{}) []byte {
	return Marshal(self, header, content)
}

func (self ProtobufEncoder) MarshalContent(obj interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(obj)
	return buffer.Bytes()
}
