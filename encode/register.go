package encode

import "github.com/jennal/goplay/pkg"

var encoderMap = map[pkg.EncodingType]EncodeDecoder{
	pkg.ENCODING_GOB:  Gob{},
	pkg.ENCODING_JSON: Json{},
	pkg.ENCODING_BSON: Bson{},
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
