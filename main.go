package main

import (
	"fmt"
	"reflect"

	"io"

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

	// cli := tcp.NewClient()
	// cli.Connect("", 8081)

	// header := pkg.NewHeader(
	// 	pkg.PKG_HEARTBEAT,
	// 	pkg.ENCODING_JSON,
	// 	"test.hello.world",
	// )
	// fmt.Println(header)
	// cli.Send(header, Message{
	// 	Id: 1,
	// 	Ok: true,
	// 	M: map[string]int{
	// 		"hello": 0,
	// 		"world": 1,
	// 	},
	// 	Arr: []string{
	// 		"from",
	// 		"client",
	// 	},
	// })
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
		header, body, err := client.Recv()
		fmt.Printf("Recv:\n%#v\n%#v\n%v\n", header, body, err)
		if err == io.EOF {
			break
		}
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

		fmt.Printf("\t")
		for j := 1; j < method.Type.NumIn(); j++ {
			inT := method.Type.In(j)
			fmt.Printf("%s ", inT.String())
		}
		fmt.Println()
	}

	// serv := tcp.NewServer("", 8081, &ServerHandler{})
	// if err := serv.Start(); err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Scanf("%s", nil)
}
