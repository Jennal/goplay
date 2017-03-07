package service

import (
	"fmt"

	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type Service struct {
	server transfer.IServer
	router router.Router
}

func NewService(serv transfer.IServer) *Service {
	instance := &Service{
		server: serv,
	}
	serv.RegistHandler(instance)
	return instance
}

func (self *Service) Regist(obj interface{}) {
	self.router.Add(obj)
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
func (self *Service) OnNewClient(client transfer.IClient) {
	fmt.Println("OnNewClient", client)
	sess := session.NewSession(client)
	go func() {
		for {
			header, bodyBuf, err := client.Recv()
			fmt.Printf("Recv:\n\theader => %#v\n\terr => %v\n", header, err)
			if err != nil {
				//TODO: log err
				break
			}

			fmt.Printf("Recv:\n\tbody => %#v\nerr => %v\n", bodyBuf, err)
			if err != nil {
				//TODO: log err
				break
			}

			switch header.Type {
			case pkg.PKG_NOTIFY:
				self.callRouteFunc(sess, header, bodyBuf)
			case pkg.PKG_REQUEST:
				results := self.callRouteFunc(sess, header, bodyBuf)
				err := self.response(sess, header, results)
				if err != nil {
					//TODO: log err
					break
				}
			case pkg.PKG_HEARTBEAT:
				//TODO:
			case pkg.PKG_HEARTBEAT_RESPONSE:
				//TODO:
			default:
				break
			}
		}
	}()
}

func (self *Service) callRouteFunc(sess *session.Session, header *pkg.Header, bodyBuf []byte) []interface{} {
	/*
	 * 1. find route func
	 * 2. unmarshal data
	 * 3. call route func
	 */
	method := self.router.Get(header.Route)
	val := method.NewArg(2)
	decoder := encode.GetEncodeDecoder(header.Encoding)
	err := decoder.Unmarshal(bodyBuf, val)
	if err != nil {
		//TODO: log err
		return nil
	}
	return method.Call(sess, val)
}

func (self *Service) response(sess *session.Session, header *pkg.Header, results []interface{}) error {
	if len(results) <= 0 {
		respHeader := *header
		return sess.Send(&respHeader, []byte{})
	}

	encoder := encode.GetEncodeDecoder(header.Encoding)
	body, err := encoder.Marshal(results[0])
	if err != nil {
		return err
	}

	respHeader := *header
	return sess.Send(&respHeader, body)
}
