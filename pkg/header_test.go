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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderDecode(t *testing.T) {
	pack := Header{
		Type:        PKG_NOTIFY,
		Encoding:    ENCODING_GOB,
		ID:          2,
		Status:      STAT_OK,
		ContentSize: 10,
		Route:       "test.header",
	}
	buffer, err := pack.Marshal()
	t.Log(buffer, string(buffer[HEADER_STATIC_SIZE+1:]), err)
	assert.Nil(t, err, "MarshalHeader error")

	newPack := Header{}
	n, err := UnmarshalHeader(buffer, &newPack)

	assert.Equal(t, n, HEADER_STATIC_SIZE+1+len(pack.Route), "Unmarshal Header size(%v) != HEADER_STATIC_SIZE+1+len(pack.Route)(%v)", n, HEADER_STATIC_SIZE+1+len(pack.Route))
	assert.Nil(t, err, "UnmarshalHeader error")

	assert.Equal(t, pack.Type, newPack.Type, "package.Type are not equal %v != %v", pack.Type, newPack.Type)
	assert.Equal(t, pack.Encoding, newPack.Encoding, "package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
	assert.Equal(t, pack.ID, newPack.ID, "package.ID are not equal %v != %v", pack.ID, newPack.ID)
	assert.Equal(t, pack.ContentSize, newPack.ContentSize, "package.ContentSize are not equal %v != %v", pack.ContentSize, newPack.ContentSize)
	assert.Equal(t, pack.Route, newPack.Route, "package.Route are not equal %v != %v", pack.Route, newPack.Route)
}

func TestHeaderHasNoRoute(t *testing.T) {
	pack := Header{
		Type:        PKG_HEARTBEAT,
		Encoding:    ENCODING_GOB,
		ID:          2,
		ContentSize: 10,
		Route:       "",
	}
	buffer, err := pack.Marshal()
	assert.Nil(t, err, "MarshalHeader error")

	newPack := Header{}
	n, err := UnmarshalHeader(buffer, &newPack)

	assert.Equal(t, n, HEADER_STATIC_SIZE+1+len(pack.Route), "Unmarshal Header size(%v) != HEADER_STATIC_SIZE+1+len(pack.Route)(%v)", n, HEADER_STATIC_SIZE+1+len(pack.Route))
	assert.Nil(t, err, "UnmarshalHeader error")

	assert.Equal(t, pack.Type, newPack.Type, "package.Type are not equal %v != %v", pack.Type, newPack.Type)
	assert.Equal(t, pack.Encoding, newPack.Encoding, "package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
	assert.Equal(t, pack.ID, newPack.ID, "package.ID are not equal %v != %v", pack.ID, newPack.ID)
	assert.Equal(t, pack.ContentSize, newPack.ContentSize, "package.ContentSize are not equal %v != %v", pack.ContentSize, newPack.ContentSize)
	assert.Equal(t, pack.Route, newPack.Route, "package.Route are not equal %v != %v", pack.Route, newPack.Route)
}

func BenchmarkHeaderDecode(b *testing.B) {
	pack := Header{
		Type:        PKG_NOTIFY,
		Encoding:    ENCODING_GOB,
		ID:          2,
		ContentSize: 3,
	}
	buffer, err := pack.Marshal()
	assert.Nil(b, err, "MarshalHeader error")

	newHeader := Header{}
	for i := 0; i < b.N; i++ {
		n, err := UnmarshalHeader(buffer, &newHeader)
		assert.Equal(b, n, HEADER_STATIC_SIZE, "Unmarshal Header size(%v) != HEADER_STATIC_SIZE(%v)", n, HEADER_STATIC_SIZE)
		assert.Nil(b, err, "UnmarshalHeader error")
	}
}
