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
	"github.com/golang/protobuf/proto"
	"github.com/jennal/goplay/log"
)

//Protobuf encoder/decoder
type Protobuf struct {
}

type oneOfMarshaler interface {
	MarshalOneOf() ([]byte, error)
	UnmarshalOneOf(buf []byte) error
}

func (self Protobuf) Marshal(obj interface{}) ([]byte, error) {
	if oneOf, ok := obj.(oneOfMarshaler); ok {
		return oneOf.MarshalOneOf()
	}

	pb, ok := obj.(proto.Message)
	if !ok {
		return nil, log.NewErrorf("protobuf: convert on wrong type value: %#v", obj)
	}
	return proto.Marshal(pb)
}

func (self Protobuf) Unmarshal(data []byte, content interface{}) error {
	if oneOf, ok := content.(oneOfMarshaler); ok {
		return oneOf.UnmarshalOneOf(data)
	}

	pb, ok := content.(proto.Message)
	if !ok {
		return log.NewErrorf("protobuf: convert on wrong type value: %#v", content)
	}
	return proto.Unmarshal(data, pb)
}

func NewProtobuf() EncodeDecoder {
	return &Base{
		child: Protobuf{},
	}
}
