package tcp

import (
	"fmt"
	"testing"

	"time"

	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
	"github.com/stretchr/testify/assert"
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
func (self *ServerHandler) OnNewClient(client transfer.IClient) {
	fmt.Println("OnNewClient", client)
	for {
		header, bodyBuf, err := client.Recv()
		fmt.Println("Recv Error: ", err)
		var obj Message
		err = encode.GetEncodeDecoder(pkg.ENCODING_JSON).Unmarshal(bodyBuf, &obj)
		fmt.Println("Recv Error: ", err)
		fmt.Printf("Recv:\nheader => %#v\nbodyBuf => %v\nmessage => %#v\n", header, bodyBuf, obj)
		if err != nil {
			break
		}
	}
}

func TestTcp(t *testing.T) {
	serv := NewServer("", 8888)
	serv.RegistHandler(&ServerHandler{})
	go serv.Start()

	cli := NewClient()
	cli.Connect("", 8888)

	header := pkg.NewHeader(
		pkg.PKG_HEARTBEAT,
		pkg.ENCODING_JSON,
		"test.hello.world",
	)
	t.Log(header)
	obj := Message{
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
	}
	buf, err := encode.GetEncodeDecoder(pkg.ENCODING_JSON).Marshal(obj)
	assert.Nil(t, err, "Encode Error: %v", err)
	cli.Send(header, buf)
	time.Sleep(time.Second)
}
