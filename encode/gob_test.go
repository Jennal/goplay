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

package encode

import (
	"fmt"
	"testing"

	"github.com/jennal/goplay/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGobDecode(t *testing.T) {
	encoder := Gob{}
	decoder := Gob{}

	content := []int{1, 2, 3}
	buffer, err := encoder.Marshal(content)
	assert.Nil(t, err, "encode.Marshal error")

	var newContent []int
	err = decoder.Unmarshal(buffer, &newContent)
	assert.Nil(t, err, "decoder.Unmarshal error")

	fmt.Println(content, newContent)
	assert.Equal(t, content[0], newContent[0], "package.Content[0] are not equal %v != %v", content[0], newContent[0])
}

type message struct {
	ID  int
	Ok  bool
	M   map[string]int
	Arr []string
}

func TestGobMarshal(t *testing.T) {
	encode := GetEncodeDecoder(pkg.ENCODING_GOB)
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

func BenchmarkGobDecode(b *testing.B) {
	encoder := Gob{}
	decoder := Gob{}
	content := []int{1, 2, 3, 4}

	buffer, err := encoder.Marshal(content)
	assert.Nil(b, err, "encode.Marshal error")
	var newContent []int
	for i := 0; i < b.N; i++ {
		err = decoder.Unmarshal(buffer, &newContent)
		assert.Nil(b, err, "decoder.Unmarshal error")
	}
}
