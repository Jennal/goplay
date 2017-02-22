package protocol

import (
	"encoding/json"

	"github.com/jennal/goplay/handler/pkg"
)

type JsonEncoder struct {
	HeaderEncoder
}

func (self JsonEncoder) Marshal(header *pkg.Header, content interface{}) []byte {
	return Marshal(self, header, content)
}

func (self JsonEncoder) MarshalContent(obj interface{}) []byte {
	result, _ := json.Marshal(obj)
	return result
}
