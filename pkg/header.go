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

const (
	HEADER_STATIC_SIZE = 6
	INT_BYTE_SIZE      = 4
)

type Header struct {
	Type        PackageType
	Encoding    EncodingType
	ID          PackageIDType
	Status      Status
	ContentSize PackageSizeType
	Route       string

	ClientID uint32 /* rpc only */
}

func NewHeader(t PackageType, e EncodingType, idGen *IDGen, r string) *Header {
	return &Header{
		Type:        t,
		Encoding:    e,
		ID:          idGen.NextID(),
		Status:      STAT_OK,
		ContentSize: 0,
		Route:       r,
	}
}

func NewRpcHeader(h *Header, clientId uint32) *Header {
	return &Header{
		Type:        h.Type | PKG_RPC,
		Encoding:    h.Encoding,
		ID:          h.ID,
		Status:      h.Status,
		ContentSize: h.ContentSize,
		Route:       h.Route,
		ClientID:    clientId,
	}
}

func NewHeaderFromRpc(h *Header) *Header {
	t := h.Type
	//For broadcast through backend
	if t == PKG_RPC_BROADCAST {
		t = PKG_PUSH
	}

	return &Header{
		Type:        t &^ PKG_RPC,
		Encoding:    h.Encoding,
		ID:          h.ID,
		Status:      h.Status,
		ContentSize: h.ContentSize,
		Route:       h.Route,
		ClientID:    0,
	}
}

func (self *Header) Marshal() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.Write([]byte{
		byte(self.Type),
		byte(self.Encoding),
		byte(self.ID),
		byte(self.Status),
	})
	buf, err := helpers.GetBytes(self.ContentSize)
	if err != nil {
		return nil, err
	}
	buffer.Write(buf)

	buffer.WriteByte(byte(len(self.Route)))
	buffer.Write([]byte(self.Route))

	if self.Type&PKG_RPC == PKG_RPC {
		buf, err := helpers.GetBytes(self.ClientID)
		if err != nil {
			return nil, err
		}
		buffer.Write(buf)
	}

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

	n, err := UnmarshalHeader(buffer, header)
	if err != nil {
		return n, err
	}

	if header.Type&PKG_RPC == PKG_RPC {
		clientIDBuf := make([]byte, INT_BYTE_SIZE)
		_, err := reader.Read(clientIDBuf)
		if err != nil {
			return 0, err
		}

		header.ClientID, err = helpers.ToUInt32(clientIDBuf)
		if err != nil {
			return 0, err
		}
	}

	return n + INT_BYTE_SIZE, nil
}

func UnmarshalHeader(data []byte, header *Header) (int, error) {
	if len(data) < HEADER_STATIC_SIZE {
		return 0, log.NewError("data size < HEADER_STATIC_SIZE")
	}

	buffer := bytes.NewBuffer(data)

	b, err := buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Type = PackageType(b)

	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Encoding = EncodingType(b)

	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.ID = PackageIDType(b)

	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Status = Status(b)

	size, err := helpers.ToUInt16(data[HEADER_STATIC_SIZE-2 : HEADER_STATIC_SIZE])
	if err != nil {
		return 0, err
	}
	header.ContentSize = PackageSizeType(size)

	for i := HEADER_STATIC_SIZE - 2; i < HEADER_STATIC_SIZE; i++ {
		buffer.ReadByte()
	}
	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	routeSize := int(b)

	if routeSize > 0 {
		route := make([]byte, b)
		_, err = buffer.Read(route)
		if err != nil {
			return 0, err
		}
		header.Route = string(route)
	}

	return HEADER_STATIC_SIZE + 1 + routeSize, nil
}
