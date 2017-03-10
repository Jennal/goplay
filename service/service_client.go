package service

import (
	"sync"
	"time"

	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

const (
	REQUEST_TIMEOUT = 3 * time.Second
)

type requestCallbacks struct {
	startTime      time.Time
	successCallbak *Method
	failCallback   *Method
}

type ServiceClient struct {
	transfer.IClient
	encoding         pkg.EncodingType
	encoder          encode.EncodeDecoder
	heartBeatManager *heartbeat.HeartBeatManager

	requestCbsMutex sync.Mutex
	requestCbs      map[pkg.PackageIDType]*requestCallbacks

	pushCbsMutex sync.Mutex
	pushCbs      map[string][]*Method
}

func NewServiceClient(cli transfer.IClient, e pkg.EncodingType) *ServiceClient {
	encoder := encode.GetEncodeDecoder(e)
	if encoder == nil {
		log.Errorf("Can't find Encoder for %v", e)
		return nil
	}

	result := &ServiceClient{
		IClient:          cli,
		encoding:         e,
		encoder:          encoder,
		heartBeatManager: heartbeat.NewHeartBeatManager(),

		requestCbs: make(map[pkg.PackageIDType]*requestCallbacks),
		pushCbs:    make(map[string][]*Method),
	}
	go result.checkTimeoutLoop()
	result.setupEventLoop()

	return result
}

func (s *ServiceClient) checkTimeoutLoop() {
	for {
		ids := []pkg.PackageIDType{}

		s.requestCbsMutex.Lock()
		for id, item := range s.requestCbs {
			if time.Since(item.startTime) > REQUEST_TIMEOUT {
				ids = append(ids, id)
				item.failCallback.Call(pkg.NewErrorMessage(pkg.STAT_ERR_TIMEOUT, "Request Timeout"))
			}
		}

		for _, id := range ids {
			delete(s.requestCbs, id)
		}
		s.requestCbsMutex.Unlock()

		time.Sleep(REQUEST_TIMEOUT)
	}
}

func (s *ServiceClient) setupEventLoop() {
	var exitChan chan int
	s.On(transfer.EVENT_CLIENT_CONNECTED, s, func(client transfer.IClient) {
		sess := session.NewSession(client)
		s.heartBeatManager.OnNewClient(sess)
		go func() {
			for {
				select {
				case <-exitChan:
					break
				default:
					header, bodyBuf, err := client.Recv()
					if header.Type != pkg.PKG_HEARTBEAT && header.Type != pkg.PKG_HEARTBEAT_RESPONSE {
						log.Logf("Recv:\n\theader => %#v\n\tbody => %#v | %v\n\terr => %v\n", header, bodyBuf, string(bodyBuf), err)
					}

					if err != nil {
						log.Errorf("Recv:\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
						client.Disconnect()
						break
					}

					switch header.Type {
					case pkg.PKG_NOTIFY:
						s.recvPush(header, bodyBuf)
					case pkg.PKG_RESPONSE:
						s.recvResponse(header, bodyBuf)
					case pkg.PKG_HEARTBEAT:
						s.heartBeatManager.OnRecv(sess, header, bodyBuf)
					case pkg.PKG_HEARTBEAT_RESPONSE:
						s.heartBeatManager.OnRecv(sess, header, bodyBuf)
					case pkg.PKG_REQUEST:
						fallthrough
					default:
						log.Errorf("Can't reach here!!\n\terr => %v\n\theader => %#v\n\tbody => %#v", err, header, bodyBuf)
						break
					}
				}
			}
		}()
	})
	s.On(transfer.EVENT_CLIENT_DISCONNECTED, s, func(cli transfer.IClient) {
		exitChan <- 1
	})
}

func (s *ServiceClient) recvPush(header *pkg.Header, body []byte) {
	s.pushCbsMutex.Lock()
	list, ok := s.pushCbs[header.Route]
	s.pushCbsMutex.Unlock()

	if !ok {
		return
	}

	for _, item := range list {
		val := item.NewArg(0)
		s.encoder.Unmarshal(body, val)
		item.Call(helpers.GetValueFromPtr(val))
	}
}

func (s *ServiceClient) recvResponse(header *pkg.Header, body []byte) {
	s.requestCbsMutex.Lock()
	cbs, ok := s.requestCbs[header.ID]
	if ok {
		delete(s.requestCbs, header.ID)
	}
	s.requestCbsMutex.Unlock()

	if !ok {
		return
	}

	val := cbs.successCallbak.NewArg(0)
	s.encoder.Unmarshal(body, val)
	cbs.successCallbak.Call(helpers.GetValueFromPtr(val))
}

func (s *ServiceClient) Request(route string, data interface{}, succCb interface{}, failCb func(*pkg.ErrorMessage)) error {
	header := s.NewHeader(pkg.PKG_REQUEST, s.encoding, route)
	cbs := requestCallbacks{
		successCallbak: NewMethod(succCb),
		failCallback:   NewMethod(failCb),
		startTime:      time.Now(),
	}

	s.requestCbsMutex.Lock()
	s.requestCbs[header.ID] = &cbs
	s.requestCbsMutex.Unlock()

	buf, err := s.encoder.Marshal(data)
	if err != nil {
		return err
	}
	return s.Send(header, buf)
}

func (s *ServiceClient) Notify(route string, data interface{}) error {
	header := s.NewHeader(pkg.PKG_NOTIFY, s.encoding, route)
	buf, err := s.encoder.Marshal(data)
	if err != nil {
		return err
	}
	return s.Send(header, buf)
}

func (s *ServiceClient) AddListener(route string, callback interface{}) {
	s.pushCbsMutex.Lock()
	defer s.pushCbsMutex.Unlock()

	list, ok := s.pushCbs[route]
	if !ok {
		list = make([]*Method, 0)
		s.pushCbs[route] = list
	}

	s.pushCbs[route] = append(s.pushCbs[route], NewMethod(callback))
}
