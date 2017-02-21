package main

import (
	"fmt"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/protocol"
)

func init() {
	fmt.Println("init-1")
}

func init() {
	fmt.Println("init-2")
}

func main() {
	// var i int32 = 1
	// buffer := new(bytes.Buffer)
	// binary.Write(buffer, binary.LittleEndian, i)
	// fmt.Println(buffer.Bytes())

	// buf := new(bytes.Buffer)
	// var num uint16 = 1234
	// err := binary.Write(buf, binary.LittleEndian, num)
	// if err != nil {
	// 	fmt.Println("binary.Write failed:", err)
	// }
	// fmt.Printf("% x", buf.Bytes())

	// buffer := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

	// fmt.Println("Hello", buffer[:2], buffer[1:9], buffer[2:])
	// buffer = append(buffer[:2], buffer[3:]...)
	// fmt.Println(buffer)

	// var buf bytes.Buffer
	// encoder := gob.NewEncoder(&buf)
	// encoder.Encode(buffer)
	// var newBuffer []byte
	// decoder := gob.NewDecoder(&buf)
	// decoder.Decode(&newBuffer)
	// fmt.Println("newBuffer", newBuffer)

	encoder := protocol.GobEncoder{}
	decoder := protocol.GobDecoder{}

	content := []int{1, 2, 3}
	pack := pkg.Header{
		Type:     pkg.PKG_NOTIFY,
		Encoding: pkg.ENCODING_GOB,
		ID:       2,
	}
	buffer := encoder.Marshal(&pack, content)

	newPack := pkg.Header{}
	var newContent []int
	decoder.Unmarshal(buffer, &newPack, &newContent)

	fmt.Println(pack, newPack)
	fmt.Println(content, newContent)
}
