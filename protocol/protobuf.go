package protocol

import (
	"bytes"
	"encoding/gob"

	"github.com/jennal/goplay/handler/pkg"
)

/* TODO: */
type Protobuf struct {
	HeaderEncoder
	HeaderDecoder
}

func (self Protobuf) Marshal(header *pkg.Header, content interface{}) []byte {
	return Marshal(self, header, content)
}

func (self Protobuf) MarshalContent(obj interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(obj)
	return buffer.Bytes()
}

func (self Protobuf) Unmarshal(data []byte, header *pkg.Header, content interface{}) {
	Unmarshal(self, data, header, content)
}

func (self Protobuf) UnmarshalContent(data []byte, content interface{}) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(content)
}
