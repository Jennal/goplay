package main

import "fmt"

func init() {
	fmt.Println("init-1")
}

func init() {
	fmt.Println("init-2")
}

func main() {
	buffer := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

	fmt.Println("Hello", buffer[:2], buffer[1:9], buffer[2:])
	buffer = append(buffer[:2], buffer[3:]...)
	fmt.Println(buffer)
}
