package main

import (
	"fmt"

	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer/tcp"
)

func init() {
	fmt.Println("init-1")
}

func init() {
	fmt.Println("init-2")
}

type Handler struct{}

func (self *Handler) OnStarted() {
	fmt.Println("Handler-OnStarted")
}

func (self *Handler) OnStopped() {
	fmt.Println("Handler-OnStopped")
}

func (self *Handler) OnNewClient(sess *session.Session) {
	fmt.Println("Handler-OnNewClient", sess)
	sess.Push("test.push", "Hello from Push")
}

func (self *Handler) Test(sess *session.Session, line string) *pkg.ErrorMessage {
	fmt.Println("Handler-Test", sess, line)
	return nil
}

func (self *Handler) Add(sess *session.Session, n int) (int, *pkg.ErrorMessage) {
	fmt.Println("Handler-Add", sess, n)
	return n + 1, nil
}

func main() {
	ser := tcp.NewServer("", 9999)
	serv := service.NewService("test", ser)

	serv.RegistHanlder(&Handler{})

	err := ser.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Scanf("%s", nil)

	cli := tcp.NewClient()
	client := service.NewServiceClient(cli, pkg.ENCODING_JSON)
	client.Connect("", 9999)

	client.Request("test.handler.add", 1, func(result int) {
		fmt.Println("[test.handler.add] callback: ", result)
	}, func(err *pkg.ErrorMessage) {
		fmt.Println("[test.handler.add] error: ", err)
	})

	// data, err := json.Marshal("Hello From Client")
	// fmt.Println("json encode:", string(data))
	// if err != nil {
	// 	fmt.Println("json.Marshal error:", err)
	// } else {
	// 	// var str string
	// 	// json.Unmarshal(data, &str)
	// 	// fmt.Println("Unmarshal:", str)
	// 	cli.Send(cli.NewHeader(pkg.PKG_NOTIFY, pkg.ENCODING_JSON, "test.handler.test"), data)
	// }

	// data, err = json.Marshal(1000)
	// fmt.Println("json encode:", string(data))
	// if err != nil {
	// 	fmt.Println("json.Marshal error:", err)
	// } else {
	// 	cli.Send(cli.NewHeader(pkg.PKG_REQUEST, pkg.ENCODING_JSON, "test.handler.add"), data)
	// }

	// go func() {
	// 	for i := 0; true; i++ {
	// 		header, body, err := cli.Recv()
	// 		fmt.Printf("[%v] Recv:\n\theader => %#v\n\tbody => %#v | %v\n\terr => %v\n", time.Now(), header, body, string(body), err)
	// 		if err != nil {
	// 			break
	// 		}
	// 		if /*i%4 == 0 &&*/ header.Type == pkg.PKG_HEARTBEAT {
	// 			respHeader := cli.NewHeartBeatResponseHeader(header)
	// 			cli.Send(respHeader, []byte{})
	// 		}
	// 	}
	// }()

	fmt.Scanf("%s", nil)
}
