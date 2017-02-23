package protocol

import (
	"github.com/jennal/goplay/handler/pkg"
	"gopkg.in/mgo.v2/bson"
)

type Bson struct {
	HeaderEncoder
	HeaderDecoder
}

func (self Bson) Marshal(header *pkg.Header, content interface{}) ([]byte, error) {
	return marshal(self, header, content)
}

func (self Bson) MarshalContent(obj interface{}) ([]byte, error) {
	return bson.Marshal(obj)
}

func (self Bson) Unmarshal(data []byte, header *pkg.Header, content interface{}) error {
	return unmarshal(self, data, header, content)
}

func (self Bson) UnmarshalContent(data []byte, content interface{}) error {
	return bson.Unmarshal(data, content)
}
