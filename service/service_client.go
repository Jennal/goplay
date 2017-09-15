// Copyright (C) 2017 Jennal(jennalcn@gmail.com). All rights reserved.
//
// Licensed under the MIT License (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package service

import (
	"sync"
	"time"

	"fmt"

	"github.com/jennal/goplay/aop"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/filter/handshake"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type requestCallbacks struct {
	startTime      time.Time
	successCallbak *Method
	failCallback   *Method
}

type ServiceClient struct {
	*session.Session
	sessionManager   *session.SessionManager
	heartBeatManager filter.IFilter

	filters []filter.IFilter

	requestCbsMutex sync.Mutex
	requestCbs      map[pkg.PackageIDType]*requestCallbacks

	pushCbsMutex sync.Mutex
	pushCbs      map[string][]*Method

	handShakeChan chan bool
}

//for client
func NewServiceClient(cli transfer.IClient) *ServiceClient {
	result := &ServiceClient{
		Session:          session.NewSession(cli),
		sessionManager:   session.NewSessionManager(),
		heartBeatManager: heartbeat.NewHeartBeatManager(),

		filters: []filter.IFilter{
		// heartbeat.NewHeartBeatManager(),
		},

		requestCbs: make(map[pkg.PackageIDType]*requestCallbacks),
		pushCbs:    make(map[string][]*Method),

		handShakeChan: make(chan bool),
	}
	// result.BindClientID(cli.Id())
	result.setupEventLoop()

	return result
}

func (self *ServiceClient) RegistFilter(filter filter.IFilter) {
	self.filters = append(self.filters, filter)
}

func (self *ServiceClient) SetFilters(filters []filter.IFilter) {
	self.filters = filters
}

func (self *ServiceClient) SetHeartBeatManager(f filter.IFilter) {
	self.heartBeatManager = f
}

func (s *ServiceClient) Connect(host string, port int) error {
	if err := s.IClient.Connect(host, port); err != nil {
		return err
	}

	// s.BindClientID(s.IClient.Id())
	sess := s.getSession(s.ID, s.ClientID)
	if s.filters != nil && len(s.filters) > 0 {
		for _, filter := range s.filters {
			if !filter.OnNewClient(sess) {
				return nil
			}
		}
	}

	<-s.handShakeChan
	return nil
}

func (s *ServiceClient) checkTimeoutLoop() {
	for {
		if !s.IsConnected() {
			break
		}

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

func (s *ServiceClient) getSession(id uint32, clientId uint32) *session.Session {
	sess := s.sessionManager.GetSessionByID(id, clientId)
	if sess == nil {
		sess = session.NewSession(s)
		sess.Bind(s.ID)
		sess.BindClientID(clientId)

		s.sessionManager.Add(sess)
	}

	return sess
}

func (s *ServiceClient) sendHandShake() {
	handshake.SendHandShake(s, s.Encoding, s.Encoder, "")
}

func (s *ServiceClient) getStringRouter(idx pkg.RouteIndex) string {
	str, ok := pkg.HandShakeInstance.GetStringRoute(idx)
	if !ok {
		return ""
	}

	return str
}

func (s *ServiceClient) setupEventLoop() {
	s.AddListener(ON_SERVICE_DOWN, func(ok bool) {
		// log.Log(ON_SERVICE_DOWN)
		s.Disconnect()
	})

	exitChan := make(chan int, 1)
	s.On(transfer.EVENT_CLIENT_CONNECTED, s, func(client transfer.IClient) {
		// sess := s.getSession(s.ID, client.Id())
		// if s.filters != nil && len(s.filters) > 0 {
		// 	for _, filter := range s.filters {
		// 		if !filter.OnNewClient(sess) {
		// 			return
		// 		}
		// 	}
		// }

		s.sendHandShake()

		//heart beat
		if !s.heartBeatManager.OnNewClient(s.Session) {
			return
		}

		go s.checkTimeoutLoop()
		go func() {
			aop.Recover(func() {
			Loop:
				for {
					select {
					case <-exitChan:
						break Loop
					default:
						header, bodyBuf, err := s.Recv()
						if err != nil {
							log.Errorf("Recv:\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
							s.Disconnect()
							break Loop
						}

						// if header.Type != pkg.PKG_HEARTBEAT && header.Type != pkg.PKG_HEARTBEAT_RESPONSE {
						// 	log.Logf("Recv:\n\theader => %#v\n\tbody => %#v | %v\n\terr => %v\n", header, bodyBuf, string(bodyBuf), err)
						// }

						clientId := header.ClientID
						if clientId == 0 {
							clientId = s.ClientID
						}

						sess := s.sessionManager.GetSessionByID(s.ID, clientId)
						if sess == nil {
							sess = session.NewSession(s)
							sess.Bind(s.ID)
							sess.BindClientID(clientId)

							s.sessionManager.Add(sess)
						}

						//filters
						if s.filters != nil && len(s.filters) > 0 {
							for _, filter := range s.filters {
								if !filter.OnRecv(sess, header, bodyBuf) {
									goto Loop
								}
							}
						}

						//heart beat
						if !s.heartBeatManager.OnRecv(s.Session, header, bodyBuf) {
							goto Loop
						}

						switch header.Type {
						case pkg.PKG_PUSH, pkg.PKG_RPC_PUSH:
							s.recvPush(header, bodyBuf)
						case pkg.PKG_RESPONSE, pkg.PKG_RPC_RESPONSE:
							s.recvResponse(header, bodyBuf)
						case pkg.PKG_HAND_SHAKE_RESPONSE:
							s.recvHandShakeResponse(header, bodyBuf)
						case pkg.PKG_HEARTBEAT, pkg.PKG_HEARTBEAT_RESPONSE:
							fallthrough
						default:
							log.Errorf("Can't reach here!!\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
							break
						}
					}
				}
			}, func(err interface{}) {
				log.RecoverErrorf("%v", err)
				s.Disconnect()
			})
		}()
	})
	s.Once(transfer.EVENT_CLIENT_DISCONNECTED, s, func(cli transfer.IClient) {
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
		s.Encoder.Unmarshal(body, val)
		// log.Log("==========>\t", i, "\t", val)
		item.Call(val)
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

	// log.Logf("%v %v %v", header.Status, body, string(body))
	if header.Status == pkg.STAT_OK {
		val := cbs.successCallbak.NewArg(0)
		err := s.Encoder.Unmarshal(body, val)
		if err == nil {
			cbs.successCallbak.Call(val)
			return
		}
	} else {
		val := cbs.failCallback.NewArg(0)
		err := s.Encoder.Unmarshal(body, val)
		if err == nil {
			cbs.failCallback.Call(val)
			return
		}
	}

	cbs.failCallback.Call(pkg.NewErrorMessage(
		pkg.STAT_ERR_DECODE_FAILED,
		fmt.Sprintf("decode body failed: %#v | %v", body, string(body))))
}

func (s *ServiceClient) recvHandShakeResponse(header *pkg.Header, body []byte) {
	log.Logf("HandShake Response:\n\theader => %#v\n\tbody => %#v | %v", header, body, string(body))
	encoder := encode.GetEncodeDecoder(header.Encoding)
	resp := &pkg.HandShakeResponse{}
	encoder.Unmarshal(body, resp)
	pkg.HandShakeInstance.UpdateHandShakeResponse(resp)

	s.handShakeChan <- true
}

func (s *ServiceClient) Request(route string, data interface{}, succCb interface{}, failCb func(*pkg.ErrorMessage)) error {
	header := s.NewHeader(pkg.PKG_RPC_REQUEST, s.Encoding, route)
	header = pkg.NewRpcHeader(header, s.ClientID)
	cbs := requestCallbacks{
		successCallbak: NewMethod(succCb),
		failCallback:   NewMethod(failCb),
		startTime:      time.Now(),
	}

	s.requestCbsMutex.Lock()
	s.requestCbs[header.ID] = &cbs
	s.requestCbsMutex.Unlock()

	buf, err := s.Encoder.Marshal(data)
	if err != nil {
		return err
	}

	return s.Send(header, buf)
}

func (s *ServiceClient) Notify(route string, data interface{}) error {
	header := s.NewHeader(pkg.PKG_RPC_NOTIFY, s.Encoding, route)
	header = pkg.NewRpcHeader(header, s.ClientID)
	buf, err := s.Encoder.Marshal(data)
	if err != nil {
		return err
	}
	return s.Send(header, buf)
}

func (s *ServiceClient) Push(route string, data interface{}) error {
	return s.Session.Push(route, data)
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
