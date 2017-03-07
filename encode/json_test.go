package encode

import (
	"fmt"
	"testing"

	"github.com/jennal/goplay/pkg"
	"github.com/stretchr/testify/assert"
)

func TestJsonDecode(t *testing.T) {
	encoder := Json{}
	decoder := Json{}

	content := []int{1, 2, 3}
	buffer, err := encoder.Marshal(content)
	assert.Nil(t, err, "encode.Marshal error")

	var newContent []int
	err = decoder.Unmarshal(buffer, &newContent)
	assert.Nil(t, err, "decoder.Unmarshal error")

	fmt.Println(content, newContent)
	assert.Equal(t, content[0], newContent[0], "package.Content[0] are not equal %v != %v", content[0], newContent[0])
}

func TestJsonMarshalContent(t *testing.T) {
	encode := GetEncodeDecoder(pkg.ENCODING_JSON)
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

func BenchmarkJsonDecode(b *testing.B) {
	encoder := Json{}
	decoder := Json{}
	content := []int{1, 2, 3, 4}

	buffer, err := encoder.Marshal(content)
	assert.Nil(b, err, "encode.Marshal error")
	var newContent []int
	for i := 0; i < b.N; i++ {
		err = decoder.Unmarshal(buffer, &newContent)
		assert.Nil(b, err, "decoder.Unmarshal error")
	}
}
