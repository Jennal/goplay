package protocol

import (
	"encoding/json"

	"github.com/jennal/goplay/handler/pkg"
)

type JsonDecoder struct {
	HeaderDecoder
}

func (self JsonDecoder) Unmarshal(data []byte, header *pkg.Header, content interface{}) {
	Unmarshal(self, data, header, content)
}

func (self JsonDecoder) UnmarshalContent(data []byte, content interface{}) {
	json.Unmarshal(data, content)
}
