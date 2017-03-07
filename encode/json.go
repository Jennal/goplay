package encode

import "encoding/json"

type Json struct {
}

func (self Json) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (self Json) Unmarshal(data []byte, content interface{}) error {
	return json.Unmarshal(data, content)
}
