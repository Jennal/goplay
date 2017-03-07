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
		ContentSize: 10,
		Route:       "test.header",
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
