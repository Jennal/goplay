// Copyright (C) 2017 Jennal(jennalcn@gmail.com). All rights reserved.
//
// Licensed under the MIT License (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software distributed 
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR 
// CONDITIONS OF ANY KIND, either express or implied. See the License for the 
// specific language governing permissions and limitations under the License.

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
