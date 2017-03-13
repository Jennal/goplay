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

package pkg

import (
	"bytes"
	"io"

	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
)

const HEADER_STATIC_SIZE = 5

type Header struct {
	Type        PackageType
	Encoding    EncodingType
	ID          PackageIDType
	ContentSize PackageSizeType
	Route       string
}

func NewHeader(t PackageType, e EncodingType, idGen *IDGen, r string) *Header {
	return &Header{
		Type:        t,
		Encoding:    e,
		ID:          idGen.NextID(),
		ContentSize: 0,
		Route:       r,
	}
}

type HeaderEncoder struct {
}

func (self *Header) Marshal() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write([]byte{
		byte(self.Type),
		byte(self.Encoding),
		byte(self.ID),
	})
	buf, err := helpers.GetBytes(self.ContentSize)
	buffer.Write(buf)
	buffer.WriteByte(byte(len(self.Route)))
	buffer.Write([]byte(self.Route))

	return buffer.Bytes(), err
}

func ReadHeader(reader io.Reader, header *Header) (int, error) {
	var buffer = make([]byte, HEADER_STATIC_SIZE)
	_, err := reader.Read(buffer)
	if err != nil {
		return 0, err
	}
	// fmt.Println("Header:", err, buffer)

	routeBuf := make([]byte, 1)
	_, err = reader.Read(routeBuf)
	if err != nil {
		return 0, err
	}
	buffer = append(buffer, routeBuf...)
	/* heartbeat/heartbeat_response has no route */
	if routeBuf[0] > 0 {
		routeBuf = make([]byte, routeBuf[0])
		_, err = reader.Read(routeBuf)
		if err != nil {
			return 0, err
		}

		buffer = append(buffer, routeBuf...)
	}

	_, err = UnmarshalHeader(buffer, header)
	// fmt.Println(header)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func UnmarshalHeader(data []byte, header *Header) (int, error) {
	if len(data) < HEADER_STATIC_SIZE {
		return -1, log.NewError("data size < HEADER_STATIC_SIZE")
	}

	buffer := bytes.NewBuffer(data)

	b, err := buffer.ReadByte()
	header.Type = PackageType(b)
	b, err = buffer.ReadByte()
	header.Encoding = EncodingType(b)
	b, err = buffer.ReadByte()
	header.ID = PackageIDType(b)

	size, err := helpers.ToUInt16(data[3:HEADER_STATIC_SIZE])
	header.ContentSize = PackageSizeType(size)
	for i := 3; i < HEADER_STATIC_SIZE; i++ {
		buffer.ReadByte()
	}
	b, err = buffer.ReadByte()
	routeSize := int(b)
	if routeSize > 0 {
		route := make([]byte, b)
		_, err = buffer.Read(route)
		header.Route = string(route)
	}

	return HEADER_STATIC_SIZE + 1 + routeSize, err
}
