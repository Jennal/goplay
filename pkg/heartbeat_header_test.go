package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeartBeat(t *testing.T) {
	pack := &Header{
		Type:         PKG_HEARTBEAT,
		Encoding:     ENCODING_GOB,
		ID:           2,
		ContentSize:  10,
		Route:        "",
		RouteEncoded: ROUTE_INDEX_NONE,
	}

	newPack := NewHeartBeatResponseHeader(pack)

	assert.Equal(t, PKG_HEARTBEAT, pack.Type)
	assert.Equal(t, PKG_HEARTBEAT_RESPONSE, newPack.Type)

	assert.Equal(t, pack.Encoding, newPack.Encoding)
	assert.Equal(t, pack.ID, newPack.ID)
	assert.Equal(t, pack.ContentSize, newPack.ContentSize)
	assert.Equal(t, pack.Route, newPack.Route)
	assert.Equal(t, pack.RouteEncoded, newPack.RouteEncoded)
}
