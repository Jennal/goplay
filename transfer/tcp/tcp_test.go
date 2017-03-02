package tcp

import (
	"fmt"
	"testing"

	"time"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/transfer"
)

type Message struct {
	Id  int
	Ok  bool
	M   map[string]int
	Arr []string
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
		fmt.Printf("Recv:\nheader => %#v\nmessage => %#v\nerr => %v\n", header, obj, err)
	}
}

func TestTcp(t *testing.T) {
	serv := NewServer("", 8888, &ServerHandler{})
	go serv.Start()

	cli := NewClient()
	cli.Connect("", 8888)

	header := pkg.NewHeader(
		pkg.PKG_HEARTBEAT,
		pkg.ENCODING_JSON,
	)
	t.Log(header)
	cli.Send(header, Message{
		Id: 1,
		Ok: true,
		M: map[string]int{
			"hello": 0,
			"world": 1,
		},
		Arr: []string{
			"from",
			"client",
		},
	})
	time.Sleep(time.Second)
}
