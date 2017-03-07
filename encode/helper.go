package encode

import (
	"bytes"

	"github.com/jennal/goplay/pkg"
)

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
	n, err := pkg.UnmarshalHeader(data, header)
	return header, n, err
}

func Marshal(header *pkg.Header, content interface{}) ([]byte, error) {
	encoder := GetEncodeDecoder(header.Encoding)
	contentBuff, err := encoder.Marshal(content)
	if err != nil {
		return nil, err
	}

	header.ContentSize = pkg.PackageSizeType(len(contentBuff))
	headerBuff, err := header.Marshal()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	buffer.Write(headerBuff)
	buffer.Write(contentBuff)

	return buffer.Bytes(), nil
}

func Unmarshal(data []byte, header *pkg.Header, content interface{}) error {
	n, err := pkg.UnmarshalHeader(data, header)
	if err != nil {
		return err
	}

	decoder := GetEncodeDecoder(header.Encoding)
	return decoder.Unmarshal(data[n:], content)
}
