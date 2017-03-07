package encode

import (
	"fmt"
	"testing"

	"github.com/jennal/goplay/pkg"
	"github.com/stretchr/testify/assert"
)

/* bson not support array */
// func TestBsonDecode(t *testing.T) {
// 	encoder := Bson{}
// 	decoder := Bson{}

// 	content := []int{1, 2, 3}
// 	pack := pkg.Header{
// 		Type:     pkg.PKG_NOTIFY,
// 		Encoding: pkg.ENCODING_BSON,
// 		ID:       2,
// 	}
// 	buffer, err := encoder.Marshal(&pack, content)
// 	assert.Nil(t, err, "encode.Marshal error")

// 	newPack := pkg.Header{}
// 	var newContent []int
// 	err = decoder.Unmarshal(buffer, &newPack, &newContent)
// 	assert.Nil(t, err, "decoder.Unmarshal error")

// 	fmt.Println(pack, newPack)
// 	fmt.Println(content, newContent)

// 	assert.Equal(t, pack.Type, newPack.Type, "package.Type are not equal %v != %v", pack.Type, newPack.Type)
// 	assert.Equal(t, pack.Encoding, newPack.Encoding, "package.Encoding are not equal %v != %v", pack.Encoding, newPack.Encoding)
// 	assert.Equal(t, pack.ID, newPack.ID, "package.ID are not equal %v != %v", pack.ID, newPack.ID)
// 	assert.Equal(t, content[0], newContent[0], "package.Content[0] are not equal %v != %v", content[0], newContent[0])
// }

func TestBsonMarshal(t *testing.T) {
	encode := GetEncodeDecoder(pkg.ENCODING_BSON)
	input := message{
		ID: 20,
		Ok: true,
		M: map[string]int{
			"haha": 100,
			"hoho": 200,
		},
		Arr: []string{
			"1",
			"2",
		},
	}
	buf, err := encode.Marshal(input)
	assert.Nil(t, err, "encode.Marshal error")
	fmt.Println(string(buf))

	var output message
	encode.Unmarshal(buf, &output)
	assert.Equal(t, input, output, "Unmarshaled Content not equals to the origin one, %#v != %#v", input, output)
}

func BenchmarkBsonDecode(b *testing.B) {
	encoder := Bson{}
	decoder := Bson{}
	content := []int{1, 2, 3, 4}

	buffer, err := encoder.Marshal(content)
	assert.Nil(b, err, "encode.Marshal error")
	var newContent []int
	for i := 0; i < b.N; i++ {
		err = decoder.Unmarshal(buffer, &newContent)
		assert.Nil(b, err, "decoder.Unmarshal error")
	}
}
