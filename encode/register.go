package encode

import "github.com/jennal/goplay/pkg"

var encoderMap = map[pkg.EncodingType]EncodeDecoder{
	pkg.ENCODING_GOB:      NewGob(),
	pkg.ENCODING_JSON:     NewJson(),
	pkg.ENCODING_BSON:     NewBson(),
	pkg.ENCODING_PROTOBUF: NewProtobuf(),
}

//Regist new EncodeDecoder
func Regist(e pkg.EncodingType, ed EncodeDecoder) {
	encoderMap[e] = ed
}

//GetEncodeDecoder gets EncodeDecoder by pkg.EncodingType
func GetEncodeDecoder(encoding pkg.EncodingType) EncodeDecoder {
	ed, ok := encoderMap[encoding]
	if ok {
		return ed
	}

	return nil
}

func GetMd5EncodeDecoder() EncodeDecoder {
	return encoderMap[pkg.ENCODING_JSON]
}
