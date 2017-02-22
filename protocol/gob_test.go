package protocol

import (
	"fmt"
	"testing"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGobDecode(t *testing.T) {
	encoder := Gob{}
	decoder := Gob{}

	content := []int{1, 2, 3}
	pack := pkg.Header{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_GOB,
		ID:       2,
	}
	buffer, err := encoder.Marshal(&pack, content)
	assert.Nil(t, err, "encode.Marshal error")

	newPack := pkg.Header{}
	var newContent []int
	err = decoder.Unmarshal(buffer, &newPack, &newContent)
	assert.Nil(t, err, "decoder.Unmarshal error")

	fmt.Println(pack, newPack)
	fmt.Println(content, newContent)

	assert.Equal(t, pack.Type, newPack.Type, "package.Type are not equal %v != %v", pack.Type, newPack.Type)
	assert.Equal(t, pack.Encoding, newPack.Encoding, "package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
	assert.Equal(t, pack.ID, newPack.ID, "package.ID are not equal %v != %v", pack.ID, newPack.ID)
	assert.Equal(t, content[0], newContent[0], "package.Content[0] are not equal %v != %v", content[0], newContent[0])
}

func BenchmarkGobDecode(b *testing.B) {
	encoder := Gob{}
	decoder := Gob{}
	pack := pkg.Header{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_GOB,
		ID:       2,
	}
	content := []int{1, 2, 3, 4}

	buffer, err := encoder.Marshal(&pack, content)
	assert.Nil(b, err, "encode.Marshal error")
	newHeader := pkg.Header{}
	var newContent []int
	for i := 0; i < b.N; i++ {
		err = decoder.Unmarshal(buffer, &newHeader, &newContent)
		assert.Nil(b, err, "decoder.Unmarshal error")
	}
}
