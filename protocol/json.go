package protocol

import (
	"encoding/json"

	"github.com/jennal/goplay/handler/pkg"
)

type Json struct {
	HeaderEncoder
	HeaderDecoder
}

func (self Json) Marshal(header *pkg.Header, content interface{}) []byte {
	return Marshal(self, header, content)
}

func (self Json) MarshalContent(obj interface{}) []byte {
	result, _ := json.Marshal(obj)
	return result
}

func (self Json) Unmarshal(data []byte, header *pkg.Header, content interface{}) {
	Unmarshal(self, data, header, content)
}

func (self Json) UnmarshalContent(data []byte, content interface{}) {
	json.Unmarshal(data, content)
}
