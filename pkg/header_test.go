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
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqual(t *testing.T, n int, buffer []byte, err error, newPack, pack Header) {
	assert.Equal(t, len(buffer), n, "Unmarshal Header size(%v) != len(buffer)(%v)", n, len(buffer))
	assert.Nil(t, err, "UnmarshalHeader error")

	assert.Equal(t, pack.Type, newPack.Type, "package.Type are not equal %v != %v", pack.Type, newPack.Type)
	assert.Equal(t, pack.Encoding, newPack.Encoding, "package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
	assert.Equal(t, pack.ID, newPack.ID, "package.ID are not equal %v != %v", pack.ID, newPack.ID)
	assert.Equal(t, pack.ContentSize, newPack.ContentSize, "package.ContentSize are not equal %v != %v", pack.ContentSize, newPack.ContentSize)
	assert.Equal(t, pack.Route, newPack.Route, "package.Route are not equal %v != %v", pack.Route, newPack.Route)
}

func TestHeaderDecode(t *testing.T) {
	HandShakeInstance.UpdateRoutesMap(map[string]RouteIndex{
		"test.header": 1,
	})

	tmpl := Header{
		Type:         PKG_PUSH,
		Encoding:     ENCODING_GOB,
		ID:           2,
		Status:       Status_OK,
		ContentSize:  10,
		Route:        "test.header",
		RouteEncoded: 1,
	}

	types := []PackageType{
		PKG_NOTIFY,
		PKG_REQUEST,
		PKG_RESPONSE,
		PKG_PUSH,
		PKG_HEARTBEAT,
		PKG_HEARTBEAT_RESPONSE,
		PKG_HAND_SHAKE,
		PKG_HAND_SHAKE_RESPONSE,
		PKG_RPC_NOTIFY,
		PKG_RPC_REQUEST,
		PKG_RPC_RESPONSE,
		PKG_RPC_PUSH,
		PKG_RPC_BROADCAST,
	}

	for _, ty := range types {
		pack := tmpl
		pack.Type = ty
		if !isNeedToEncodeRoute(&pack) && ty&^PKG_RPC != PKG_PUSH {
			pack.Route = ""
			pack.RouteEncoded = ROUTE_INDEX_NONE
		}
		buffer, err := pack.Marshal()
		t.Logf("%v => %v | %v | %v", ty, pack, buffer, err)
		assert.Nil(t, err, "MarshalHeader error")

		newPack := Header{}
		n, err := UnmarshalHeader(buffer, &newPack)

		assertEqual(t, n, buffer, err, newPack, pack)
	}
}

func TestHeaderReader(t *testing.T) {
	HandShakeInstance.UpdateRoutesMap(map[string]RouteIndex{
		"test.header": 1,
	})

	tmpl := Header{
		Type:         PKG_PUSH,
		Encoding:     ENCODING_GOB,
		ID:           2,
		Status:       Status_OK,
		ContentSize:  10,
		Route:        "test.header",
		RouteEncoded: 1,
	}

	types := []PackageType{
		PKG_NOTIFY,
		PKG_REQUEST,
		PKG_RESPONSE,
		PKG_PUSH,
		PKG_HEARTBEAT,
		PKG_HEARTBEAT_RESPONSE,
		PKG_HAND_SHAKE,
		PKG_HAND_SHAKE_RESPONSE,
		PKG_RPC_NOTIFY,
		PKG_RPC_REQUEST,
		PKG_RPC_RESPONSE,
		PKG_RPC_PUSH,
		PKG_RPC_BROADCAST,
	}

	for _, ty := range types {
		pack := tmpl
		if !isNeedToEncodeRoute(&pack) && ty&^PKG_RPC != PKG_PUSH {
			pack.Route = ""
			pack.RouteEncoded = ROUTE_INDEX_NONE
		}

		pack.Type = ty
		buffer, err := pack.Marshal()
		t.Logf("%#v | %v", buffer, err)
		assert.Nil(t, err, "MarshalHeader error")

		newPack := Header{}
		buf := bytes.NewBuffer(buffer)
		n, err := ReadHeader(buf, &newPack)

		assertEqual(t, n, buffer, err, newPack, pack)
		// if i == 1 {
		// 	break
		// }
	}
}

func TestHeaderHasNoRoute(t *testing.T) {
	tmpl := Header{
		Type:         PKG_HEARTBEAT,
		Encoding:     ENCODING_GOB,
		ID:           2,
		ContentSize:  10,
		Route:        "",
		RouteEncoded: ROUTE_INDEX_NONE,
	}

	types := []PackageType{
		PKG_NOTIFY,
		PKG_REQUEST,
		PKG_RESPONSE,
		PKG_PUSH,
		PKG_HEARTBEAT,
		PKG_HEARTBEAT_RESPONSE,
		PKG_HAND_SHAKE,
		PKG_HAND_SHAKE_RESPONSE,
		PKG_RPC_NOTIFY,
		PKG_RPC_REQUEST,
		PKG_RPC_RESPONSE,
		PKG_RPC_PUSH,
		PKG_RPC_BROADCAST,
	}

	for _, ty := range types {
		pack := tmpl
		if !isNeedToEncodeRoute(&pack) && ty&^PKG_RPC != PKG_PUSH {
			pack.Route = ""
			pack.RouteEncoded = ROUTE_INDEX_NONE
		}

		pack.Type = ty
		buffer, err := pack.Marshal()
		t.Logf("%#v | %v", buffer, err)
		assert.Nil(t, err, "MarshalHeader error")

		newPack := Header{}
		n, err := UnmarshalHeader(buffer, &newPack)

		assertEqual(t, n, buffer, err, newPack, pack)
	}
}

func BenchmarkHeaderDecode(b *testing.B) {
	pack := Header{
		Type:         PKG_NOTIFY,
		Encoding:     ENCODING_GOB,
		ID:           2,
		ContentSize:  3,
		RouteEncoded: ROUTE_INDEX_NONE,
	}
	buffer, err := pack.Marshal()
	assert.Nil(b, err, "MarshalHeader error")

	newHeader := Header{}
	for i := 0; i < b.N; i++ {
		//n, err := UnmarshalHeader(buffer, &newHeader)
		UnmarshalHeader(buffer, &newHeader)
		// assertEqual(, n, buffer, err, newHeader, pack)
	}
}
