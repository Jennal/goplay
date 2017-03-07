package service

import (
	"fmt"

	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/transfer"
)

type Service struct {
	server transfer.Server
	router router.Router
}

func NewService(serv transfer.Server) *Service {
	instance := &Service{
		server: serv,
	}
	serv.SetupHandler(instance)
	return instance
}

func (self *Service) OnStarted() {
	fmt.Printf("OnStarted %p\n", self)
}
func (self *Service) OnError(err error) {
	fmt.Println("OnError", err)
}
func (self *Service) OnStopped() {
	fmt.Println("OnStopped")
}
func (self *Service) OnNewClient(client transfer.Client) {
	fmt.Println("OnNewClient", client)
	for {
		header, bodyBuf, err := client.Recv()
		fmt.Printf("Recv:\n\theader => %#v\n\terr => %v\n", header, err)
		if err != nil {
			break
		}

		fmt.Printf("Recv:\n\tbody => %#v\nerr => %v\n", bodyBuf, err)
		if err != nil {
			break
		}

		/*
		 * 1. find route func
		 * 2. unmarshal data
		 * 3. call route func
		 */
		//TODO:
		// method := self.router.Get(header.Route)
	}
}
