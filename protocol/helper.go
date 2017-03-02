package protocol

import (
	"bytes"

	"github.com/jennal/goplay/handler/pkg"
)

var headerDecoder HeaderDecoder

var encoder_map map[pkg.EncodingType]EncodeDecoder = map[pkg.EncodingType]EncodeDecoder{
	pkg.ENCODING_GOB:  Gob{},
	pkg.ENCODING_JSON: Json{},
	pkg.ENCODING_BSON: Bson{},
}

func GetEncodeDecoder(encoding pkg.EncodingType) EncodeDecoder {
	return encoder_map[encoding]
}

func UnMarshalHeader(data []byte) (*pkg.Header, int, error) {
	header := &pkg.Header{}
	n, err := headerDecoder.UnmarshalHeader(data[:HEADER_SIZE], header)
	return header, n, err
}

func marshal(self Encoder, header *pkg.Header, content interface{}) ([]byte, error) {
	contentBuff, err := self.MarshalContent(content)
	if err != nil {
		return nil, err
	}

	header.ContentSize = pkg.PackageSizeType(len(contentBuff))
	headerBuff, err := self.MarshalHeader(header)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.Write(headerBuff)
	buffer.Write(contentBuff)

	return buffer.Bytes(), nil
}

func unmarshal(self Decoder, data []byte, header *pkg.Header, content interface{}) error {
	n, err := self.UnmarshalHeader(data, header)
	if err != nil {
		return err
	}
	return self.UnmarshalContent(data[n:], content)
}
