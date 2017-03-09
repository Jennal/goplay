package main

import (
	"fmt"

	"time"

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
}

func (self *Handler) Test(sess *session.Session, line string) error {
	fmt.Println("Handler-Test", sess, line)
	return nil
}

func main() {
	ser := tcp.NewServer("", 9990)
	serv := service.NewService("test", ser)

	serv.RegistHanlder(&Handler{})

	err := ser.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	cli := tcp.NewClient()
	cli.Connect("", 9990)

	for {
		header, body, err := cli.Recv()
		fmt.Println(time.Now(), header, body, err)
		if header.Type == pkg.PKG_HEARTBEAT {
			respHeader := cli.NewHeartBeatResponseHeader(header)
			cli.Send(respHeader, []byte{})
		}
	}

	fmt.Scanf("%s", nil)
}
