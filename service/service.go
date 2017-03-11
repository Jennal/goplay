package service

import (
	"fmt"

	"reflect"

	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type Service struct {
	transfer.IServer
	router *router.Router

	Name     string
	Encoding pkg.EncodingType

	handlers []handler.IHandler
	filters  []filter.IFilter
}

func NewService(name string, serv transfer.IServer) *Service {
	instance := &Service{
		Name:     name,
		Encoding: defaults.Encoding,
		IServer:  serv,
		router:   router.NewRouter(name),
	}

	serv.RegistDelegate(instance)
	instance.RegistFilter(heartbeat.NewHeartBeatManager())

	return instance
}

func (self *Service) SetEncoding(e pkg.EncodingType) error {
	if encoder := encode.GetEncodeDecoder(e); encoder != nil {
		self.Encoding = e
		return nil
	}

	return log.NewErrorf("can't find encoder with: %v", e)
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
	sess.SetEncoding(self.Encoding)

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
			if header.Type != pkg.PKG_HEARTBEAT && header.Type != pkg.PKG_HEARTBEAT_RESPONSE {
				log.Logf("Recv:\n\theader => %#v\n\tbody => %#v | %v\n\terr => %v\n", header, bodyBuf, string(bodyBuf), err)
			}

			if err != nil {
				log.Errorf("Recv:\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
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
				// fmt.Printf(" => Loop result: %#v\n", results)
				err := self.response(sess, header, results)
				if err != nil {
					log.Errorf("Response:\n\terr => %v\n\theader => %#v\n\tresults => %#v", err, header, results)
					break
				}
			case pkg.PKG_HEARTBEAT: /* Can not come to here */
				fallthrough
			case pkg.PKG_HEARTBEAT_RESPONSE: /* Can not come to here */
				fallthrough
			default:
				log.Errorf("Can't reach here!!\n\terr => %v\n\theader => %#v\n\tbody => %#v", err, header, bodyBuf)
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
		log.Errorf("Can't find method with route: %s", header.Route)
		return nil
	}
	val := method.NewArg(2)
	// fmt.Printf("Service.callRouteFunc: %#v => %v\n", val, reflect.TypeOf(val))
	decoder := encode.GetEncodeDecoder(header.Encoding)
	err := decoder.Unmarshal(bodyBuf, val)
	if err != nil {
		log.Errorf("Service.callRouteFunc decoder.Unmarshal failed: %v", err)
		return nil
	}
	// fmt.Printf("Service.callRouteFunc: %#v => %v\n", val, reflect.TypeOf(val))
	return method.Call(sess, helpers.GetValueFromPtr(val))
}

func (self *Service) response(sess *session.Session, header *pkg.Header, results []interface{}) error {
	respHeader := *header
	respHeader.Type = pkg.PKG_RESPONSE

	if results == nil || len(results) <= 0 {
		return sess.Send(&respHeader, []byte{})
	}

	result := results[0]
	/* check error != nil */
	if len(results) == 2 && !reflect.ValueOf(results[1]).IsNil() {
		result = results[1]
	}

	// fmt.Println("result:", result)

	encoder := encode.GetEncodeDecoder(header.Encoding)
	body, err := encoder.Marshal(result)
	if err != nil {
		return err
	}

	return sess.Send(&respHeader, body)
}
