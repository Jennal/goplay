package encode

import "gopkg.in/mgo.v2/bson"

type Bson struct {
}

func (self Bson) Marshal(obj interface{}) ([]byte, error) {
	return bson.Marshal(obj)
}

func (self Bson) Unmarshal(data []byte, content interface{}) error {
	return bson.Unmarshal(data, content)
}
