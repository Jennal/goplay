package encode

import (
	"bytes"
	"encoding/gob"
)

type Gob struct {
}

func (self Gob) Marshal(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(obj)
	return buffer.Bytes(), err
}

func (self Gob) Unmarshal(data []byte, content interface{}) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(content)
}
