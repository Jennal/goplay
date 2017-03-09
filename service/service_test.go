package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/service"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer/tcp"
)

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

func TestService(t *testing.T) {
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

	data, err := json.Marshal("Hello")
	if err != nil {
		fmt.Println("json.Marshal error:", err)
	} else {
		var str string
		json.Unmarshal(data, &str)
		fmt.Println("Unmarshal:", str)
		cli.Send(cli.NewHeader(pkg.PKG_NOTIFY, pkg.ENCODING_JSON, "test.handler.test"), data)
	}

	go func() {
		for i := 0; true; i++ {
			header, _, err := cli.Recv()
			// fmt.Println(time.Now(), header, body, err)
			if err != nil {
				break
			}
			if /*i%4 == 0 &&*/ header.Type == pkg.PKG_HEARTBEAT {
				respHeader := cli.NewHeartBeatResponseHeader(header)
				cli.Send(respHeader, []byte{})
			}
		}
	}()
}
