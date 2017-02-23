package main

import (
	"fmt"

	"reflect"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/transfer"
)

func init() {
	fmt.Println("init-1")
}

func init() {
	fmt.Println("init-2")
}

type ServerHandler struct {
}

func (self *ServerHandler) OnStarted() {
	fmt.Printf("OnStarted %p\n", self)
}
func (self *ServerHandler) OnError(err error) {
	fmt.Println("OnError", err)
}
func (self *ServerHandler) OnStopped() {
	fmt.Println("OnStopped")
}
func (self *ServerHandler) OnNewClient(client transfer.Client) {
	fmt.Println("OnNewClient", client)
	for {
		header := &pkg.Header{}
		var obj Message
		err := client.Recv(header, &obj)
		fmt.Printf("Recv:\n%#v\n%#v\n%v\n", header, obj, err)
	}
}

func (self *ServerHandler) onNewClient(client transfer.Client) {
	fmt.Println("OnNewClient", client)
	for {
		header := &pkg.Header{}
		var obj Message
		err := client.Recv(header, &obj)
		fmt.Printf("Recv:\n%#v\n%#v\n%v\n", header, obj, err)
	}
}

type Message struct {
	Id  int
	Ok  bool
	M   map[string]int
	Arr []string
}

func main() {
	var inst interface{} = &ServerHandler{}
	val := reflect.TypeOf(inst)
	fmt.Println("name:", val.String())
	for i := 0; i < val.NumMethod(); i++ {
		method := val.Method(i)
		fmt.Println("method:", method.Name, method.Type.String())

		if method.Name == "OnStarted" {
			fmt.Printf("%p\n", inst)
			method.Func.Call([]reflect.Value{
				reflect.ValueOf(inst),
			})
		}
	}

	// fmt.Scanf("%s", nil)
}
