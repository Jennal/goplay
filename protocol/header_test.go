package protocol

import (
	"testing"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/stretchr/testify/assert"
)

func TestHeaderDecode(t *testing.T) {
	encoder := HeaderEncoder{}
	decoder := HeaderDecoder{}
	pack := pkg.Header{
		Type:        pkg.PKG_NOTIFY,
		Encoding:    pkg.ENCODING_GOB,
		ID:          2,
		ContentSize: 10,
	}
	buffer, err := encoder.MarshalHeader(&pack)
	assert.Nil(t, err, "encoder.MarshalHeader error")

	newPack := pkg.Header{}
	n, err := decoder.UnmarshalHeader(buffer, &newPack)

	assert.Equal(t, n, HEADER_SIZE, "Unmarshal Header size(%v) != HEADER_SIZE(%v)", n, HEADER_SIZE)
	assert.Nil(t, err, "decoder.UnmarshalHeader error")

	assert.Equal(t, pack.Type, newPack.Type, "package.Type are not equal %v != %v", pack.Type, newPack.Type)
	assert.Equal(t, pack.Encoding, newPack.Encoding, "package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
	assert.Equal(t, pack.ID, newPack.ID, "package.ID are not equal %v != %v", pack.ID, newPack.ID)
	assert.Equal(t, pack.ContentSize, newPack.ContentSize, "package.ContentSize are not equal %v != %v", pack.ContentSize, newPack.ContentSize)
}

func BenchmarkHeaderDecode(b *testing.B) {
	encoder := HeaderEncoder{}
	decoder := HeaderDecoder{}
	pack := pkg.Header{
		Type:        pkg.PKG_NOTIFY,
		Encoding:    pkg.ENCODING_GOB,
		ID:          2,
		ContentSize: 3,
	}
	buffer, err := encoder.MarshalHeader(&pack)
	assert.Nil(b, err, "encoder.MarshalHeader error")

	newHeader := pkg.Header{}
	for i := 0; i < b.N; i++ {
		n, err := decoder.UnmarshalHeader(buffer, &newHeader)
		assert.Equal(b, n, HEADER_SIZE, "Unmarshal Header size(%v) != HEADER_SIZE(%v)", n, HEADER_SIZE)
		assert.Nil(b, err, "decoder.UnmarshalHeader error")
	}
}
