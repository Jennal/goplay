package protocol

import (
	"encoding/json"

	"github.com/jennal/goplay/handler/pkg"
)

type Json struct {
	HeaderEncoder
	HeaderDecoder
}

func (self Json) Marshal(header *pkg.Header, content interface{}) ([]byte, error) {
	return Marshal(self, header, content)
}

func (self Json) MarshalContent(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (self Json) Unmarshal(data []byte, header *pkg.Header, content interface{}) error {
	return Unmarshal(self, data, header, content)
}

func (self Json) UnmarshalContent(data []byte, content interface{}) error {
	return json.Unmarshal(data, content)
}
