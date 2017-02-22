package protocol

import (
	"bytes"
	"encoding/gob"

	"github.com/jennal/goplay/handler/pkg"
)

type Gob struct {
	HeaderEncoder
	HeaderDecoder
}

func (self Gob) Marshal(header *pkg.Header, content interface{}) ([]byte, error) {
	return Marshal(self, header, content)
}

func (self Gob) MarshalContent(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(obj)
	return buffer.Bytes(), err
}

func (self Gob) Unmarshal(data []byte, header *pkg.Header, content interface{}) error {
	return Unmarshal(self, data, header, content)
}

func (self Gob) UnmarshalContent(data []byte, content interface{}) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(content)
}
