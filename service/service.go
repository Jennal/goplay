package service

import (
	"fmt"

	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type Service struct {
	name   string
	server transfer.IServer
	router *router.Router

	handlers []handler.IHandler
	filters  []filter.IFilter
}

func NewService(name string, serv transfer.IServer) *Service {
	instance := &Service{
		name:   name,
		server: serv,
		router: router.NewRouter(name),
	}

	serv.RegistDelegate(instance)
	instance.RegistFilter(heartbeat.NewHeartBeatManager())

	return instance
}

func (self *Service) RegistHanlder(obj handler.IHandler) {
	self.router.Add(obj)
	self.handlers = append(self.handlers, obj)
}

func (self *Service) RegistFilter(obj filter.IFilter) {
	self.filters = append(self.filters, obj)
}

func (self *Service) OnStarted() {
	fmt.Printf("OnStarted %p\n", self)
	for _, handler := range self.handlers {
		handler.OnStarted()
	}
}
func (self *Service) OnError(err error) {
	fmt.Println("OnError", err)
}
func (self *Service) OnStopped() {
	fmt.Println("OnStopped")
	for _, handler := range self.handlers {
		handler.OnStopped()
	}
}
func (self *Service) OnNewClient(client transfer.IClient) {
	fmt.Println("OnNewClient", client)
	sess := session.NewSession(client)

	for _, filter := range self.filters {
		if !filter.OnNewClient(sess) {
			return
		}
	}

	for _, handler := range self.handlers {
		handler.OnNewClient(sess)
	}

	go func() {
		for {
		NextLoop:
			header, bodyBuf, err := client.Recv()
			fmt.Printf("Recv:\n\theader => %#v\n\tbody => %#v\n\terr => %v\n", header, bodyBuf, err)
			if err != nil {
				//TODO: log err
				break
			}

			//filters
			for _, filter := range self.filters {
				if !filter.OnRecv(sess, header, bodyBuf) {
					goto NextLoop
				}
			}

			//map to handler
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
			case pkg.PKG_HEARTBEAT: /* Can not come to here */
				fallthrough
			case pkg.PKG_HEARTBEAT_RESPONSE: /* Can not come to here */
				fallthrough
			default:
				//TODO: log err
				fmt.Printf("What?? Can't Reach Here!! Recv:\n\theader => %#v\n\tbody => %#v\n\terr => %v\n", header, bodyBuf, err)
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
	if method == nil {
		//TODO: log err
		fmt.Println("Service.callRouteFunc error: can't find method with route", header.Route)
		return nil
	}
	val := method.NewArg(2)
	decoder := encode.GetEncodeDecoder(header.Encoding)
	err := decoder.Unmarshal(bodyBuf, &val)
	if err != nil {
		//TODO: log err
		fmt.Println("Service.callRouteFunc error: decoder.Unmarshal failed", err)
		return nil
	}
	return method.Call(sess, val)
}

func (self *Service) response(sess *session.Session, header *pkg.Header, results []interface{}) error {
	if results == nil || len(results) <= 0 {
		respHeader := *header
		return sess.Send(&respHeader, []byte{})
	}

	result := results[0]
	/* check error != nil */
	if len(results) == 2 && results[1] != nil {
		result = results[1]
	}

	encoder := encode.GetEncodeDecoder(header.Encoding)
	body, err := encoder.Marshal(result)
	if err != nil {
		return err
	}

	respHeader := *header
	return sess.Send(&respHeader, body)
}
