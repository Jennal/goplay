// Code generated by "stringer -type=EncodingType"; DO NOT EDIT.

package pkg

import "fmt"

const _EncodingType_name = "ENCODING_NONEENCODING_GOBENCODING_JSONENCODING_BSONENCODING_PROTOBUF"

var _EncodingType_index = [...]uint8{0, 13, 25, 38, 51, 68}

func (i EncodingType) String() string {
	if i >= EncodingType(len(_EncodingType_index)-1) {
		return fmt.Sprintf("EncodingType(%d)", i)
	}
	return _EncodingType_name[_EncodingType_index[i]:_EncodingType_index[i+1]]
}