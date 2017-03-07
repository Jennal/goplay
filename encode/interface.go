package encode

type Encoder interface {
	Marshal(content interface{}) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, content interface{}) error
}

type EncodeDecoder interface {
	Encoder
	Decoder
}
