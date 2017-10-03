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
	"io"

	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
)

const (
	HEADER_STATIC_SIZE = 6
	INT_BYTE_SIZE      = 4
)

type Header struct {
	Type         PackageType
	Encoding     EncodingType
	ID           PackageIDType
	Status       Status
	ContentSize  PackageSizeType
	Route        string
	RouteEncoded RouteIndex

	ClientID uint32 /* rpc only */
}

func NewHeader(t PackageType, e EncodingType, idGen *IDGen, r string) *Header {
	header := &Header{
		Type:         t,
		Encoding:     e,
		ID:           idGen.NextID(),
		Status:       Status_OK,
		ContentSize:  0,
		Route:        r,
		RouteEncoded: ROUTE_INDEX_NONE,
	}

	fillIndexRoute(header)
	return header
}

func NewRpcHeader(h *Header, clientId uint32) *Header {
	if h.Type&PKG_RPC != PKG_PUSH {
		if encoded, ok := DefaultHandShake().ConvertRouteIndexToRpc(h.Route); ok {
			h.RouteEncoded = encoded
		}
	}

	return &Header{
		Type:         h.Type | PKG_RPC,
		Encoding:     h.Encoding,
		ID:           h.ID,
		Status:       h.Status,
		ContentSize:  h.ContentSize,
		Route:        h.Route,
		RouteEncoded: h.RouteEncoded,
		ClientID:     clientId,
	}
}

func NewHeaderFromRpc(h *Header) *Header {
	t := h.Type
	//For broadcast through backend
	if t == PKG_RPC_BROADCAST {
		t = PKG_PUSH
	}

	if t&^PKG_RPC != PKG_PUSH {
		if encoded, ok := DefaultHandShake().ConvertRouteIndexFromRpc(h.Route); ok {
			h.RouteEncoded = encoded
		}
	}

	return &Header{
		Type:         t &^ PKG_RPC,
		Encoding:     h.Encoding,
		ID:           h.ID,
		Status:       h.Status,
		ContentSize:  h.ContentSize,
		Route:        h.Route,
		RouteEncoded: h.RouteEncoded,
		ClientID:     0,
	}
}

func (self *Header) Marshal() ([]byte, error) {
	var buffer helpers.Buffer

	buffer.Write([]byte{
		byte(self.Type),
		byte(self.Encoding),
		byte(self.ID),
		byte(self.Status),
	})
	buffer.WriteUInt16(self.ContentSize)

	if self.Type&^PKG_RPC == PKG_PUSH {
		//write original route
		buffer.WriteByte(byte(len(self.Route)))
		buffer.Write([]byte(self.Route))
	} else {
		//write encoded route
		buffer.WriteUInt16(self.RouteEncoded)
	}

	if self.Type.IsRPC() {
		buffer.WriteUInt32(self.ClientID)
	}

	return buffer.Bytes(), nil
}

func ReadHeader(reader io.Reader, header *Header) (int, error) {
	var buffer = make([]byte, HEADER_STATIC_SIZE)
	_, err := reader.Read(buffer)
	if err != nil {
		return 0, err
	}
	// fmt.Println("Header:", err, buffer)
	n, err := UnmarshalStaticHeader(buffer, header)
	if err != nil {
		return n, err
	}
	// log.Logf("header = %#v", header)

	buffer = []byte{}
	if header.Type&^PKG_RPC == PKG_PUSH {
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
	} else {
		routeBuf := make([]byte, 2)
		_, err = reader.Read(routeBuf)
		if err != nil {
			return 0, err
		}
		buffer = append(buffer, routeBuf...)
	}

	m, err := UnmarshalDynamicHeader(buffer, header)
	if err != nil {
		return n, err
	}

	if header.Type.IsRPC() {
		clientIDBuf := make([]byte, INT_BYTE_SIZE)
		_, err := reader.Read(clientIDBuf)
		if err != nil {
			return 0, err
		}

		header.ClientID, err = helpers.ToUInt32(clientIDBuf)
		if err != nil {
			return 0, err
		}

		return n + m + INT_BYTE_SIZE, nil
	}

	return n + m, nil
}

func UnmarshalHeader(buffer []byte, header *Header) (int, error) {
	n, err := UnmarshalStaticHeader(buffer, header)
	if err != nil {
		return n, err
	}
	// log.Logf("n = %v", n)

	buffer = buffer[n:]
	m, err := UnmarshalDynamicHeader(buffer, header)
	if err != nil {
		return n + m, err
	}
	// log.Logf("m = %v", m)

	buffer = buffer[m:]
	if header.Type.IsRPC() {
		header.ClientID, err = helpers.ToUInt32(buffer[:4])
		if err != nil {
			return 0, err
		}

		return n + m + INT_BYTE_SIZE, nil
	}

	return n + m, nil
}

func UnmarshalStaticHeader(data []byte, header *Header) (int, error) {
	if len(data) < HEADER_STATIC_SIZE {
		return 0, log.NewError("data size < HEADER_STATIC_SIZE")
	}

	readSize := 0
	buffer := helpers.NewBuffer(data)

	b, err := buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Type = PackageType(b)
	readSize += 1

	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Encoding = EncodingType(b)
	readSize += 1

	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.ID = PackageIDType(b)
	readSize += 1

	b, err = buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	header.Status = Status(b)
	readSize += 1

	size, err := buffer.ReadUInt16()
	if err != nil {
		return 0, err
	}
	header.ContentSize = PackageSizeType(size)
	readSize += 2

	return readSize, nil
}

func UnmarshalDynamicHeader(data []byte, header *Header) (int, error) {
	readSize := 0
	buffer := helpers.NewBuffer(data)

	if header.Type&^PKG_RPC == PKG_PUSH {
		b, err := buffer.ReadByte()
		if err != nil {
			return 0, err
		}
		routeSize := int(b)
		readSize += 1

		if routeSize > 0 {
			route := make([]byte, b)
			_, err = buffer.Read(route)
			if err != nil {
				return 0, err
			}
			header.Route = string(route)
			readSize += len(route)

			header.RouteEncoded = ROUTE_INDEX_NONE
		}
	} else {
		r, err := buffer.ReadUInt16()
		if err != nil {
			return 0, err
		}
		header.RouteEncoded = RouteIndex(r)
		readSize += 2

		fillStringRoute(header)
		// log.Logf("index=%v, route=%v", header.RouteEncoded, header.Route)
	}

	return readSize, nil
}

func isNeedToEncodeRoute(header *Header) bool {
	return header.Type&^PKG_RPC == PKG_REQUEST ||
		header.Type&^PKG_RPC == PKG_NOTIFY ||
		header.Type&^PKG_RPC == PKG_RESPONSE
}

func fillIndexRoute(header *Header) {
	if !isNeedToEncodeRoute(header) {
		return
	}

	header.RouteEncoded, _ = DefaultHandShake().GetIndexRoute(header.Route)
}

func fillStringRoute(header *Header) {
	if !isNeedToEncodeRoute(header) {
		return
	}

	header.Route, _ = DefaultHandShake().GetStringRoute(header.RouteEncoded)
}
