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
	"reflect"
	"sync"
	"time"

	"fmt"

	"github.com/jennal/goplay/aop"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type ProcessorClient struct {
	*session.Session
	sessionManager   *session.SessionManager
	heartBeatManager filter.IFilter

	router  *router.Router
	filters []filter.IFilter

	requestCbsMutex sync.Mutex
	requestCbs      map[pkg.PackageIDType]*requestCallbacks

	pushCbsMutex sync.Mutex
	pushCbs      map[string][]*Method
}

//for server
func NewProcessorClient(cli transfer.IClient) *ProcessorClient {
	result := &ProcessorClient{
		Session:          session.NewSession(cli),
		sessionManager:   session.NewSessionManager(),
		heartBeatManager: heartbeat.NewHeartBeatManager(),

		router:  nil,
		filters: []filter.IFilter{
		// heartbeat.NewHeartBeatManager(),
		},

		requestCbs: make(map[pkg.PackageIDType]*requestCallbacks),
		pushCbs:    make(map[string][]*Method),
	}
	// result.BindClientID(cli.Id())
	result.setupEventLoop()

	return result
}

func (self *ProcessorClient) SetRouter(router *router.Router) {
	self.router = router
}

func (self *ProcessorClient) RegistFilter(filter filter.IFilter) {
	self.filters = append(self.filters, filter)
}

func (self *ProcessorClient) SetFilters(filters []filter.IFilter) {
	self.filters = filters
}

func (self *ProcessorClient) SetHeartBeatManager(f filter.IFilter) {
	self.heartBeatManager = f
}

func (s *ProcessorClient) Connect(host string, port int) error {
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

	return nil
}

func (s *ProcessorClient) checkTimeoutLoop() {
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

func (s *ProcessorClient) getSession(id uint32, clientId uint32) *session.Session {
	sess := s.sessionManager.GetSessionByID(id, clientId)
	if sess == nil {
		sess = session.NewSession(s)
		sess.Bind(s.ID)
		sess.BindClientID(clientId)

		s.sessionManager.Add(sess)
	}

	return sess
}

func (s *ProcessorClient) getStringRouter(idx pkg.RouteIndex) string {
	str, ok := pkg.HandShakeInstance.GetStringRoute(idx)
	if !ok {
		return ""
	}

	return str
}

func (s *ProcessorClient) setupEventLoop() {
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
						case pkg.PKG_REQUEST, pkg.PKG_RPC_REQUEST:
							if s.router != nil {
								results, err := s.callRouteFunc(sess, header, bodyBuf)
								if err != nil {
									log.Errorf("CallRouteFunc:\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
									sess.Disconnect()
									break Loop
								}
								// fmt.Printf(" => Loop result: %#v\n", results)
								err = s.response(sess, header, results)
								if err != nil {
									log.Errorf("Response:\n\terr => %v\n\theader => %#v\n\tresults => %#v", err, header, results)
									sess.Disconnect()
									break Loop
								}
							}
						case pkg.PKG_NOTIFY, pkg.PKG_RPC_NOTIFY:
							if s.router != nil {
								_, err := s.callRouteFunc(sess, header, bodyBuf)
								if err != nil {
									log.Errorf("CallRouteFunc:\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
									sess.Disconnect()
									break Loop
								}
							}
						case pkg.PKG_PUSH, pkg.PKG_RPC_PUSH:
							s.recvPush(header, bodyBuf)
						case pkg.PKG_RESPONSE, pkg.PKG_RPC_RESPONSE:
							s.recvResponse(header, bodyBuf)
						case pkg.PKG_HEARTBEAT, pkg.PKG_HEARTBEAT_RESPONSE:
							fallthrough
						default:
							log.Errorf("Can't reach here!!\n\terr => %v\n\theader => %#v\n\tbody => %#v | %v", err, header, bodyBuf, string(bodyBuf))
							break
						}
					}
				}
			}, func(err interface{}) {
				if e, ok := err.(error); ok {
					log.Error(e)
				} else if err != nil {
					log.Errorf("%#v", err)
				} else {
					log.Errorf("%#v", err)
				}

				s.Disconnect()
			})
		}()
	})
	s.Once(transfer.EVENT_CLIENT_DISCONNECTED, s, func(cli transfer.IClient) {
		exitChan <- 1
	})
}

func (s *ProcessorClient) callRouteFunc(sess *session.Session, header *pkg.Header, bodyBuf []byte) ([]interface{}, error) {
	/*
	 * 1. find route func
	 * 2. unmarshal data
	 * 3. call route func
	 */
	method := s.router.Get(header.Route)
	if method == nil {
		return nil, log.NewErrorf("Can't find method with route: %s", header.Route)
	}
	val := method.NewArg(2)
	// fmt.Printf("Service.callRouteFunc: %#v => %v\n", val, reflect.TypeOf(val))
	decoder := encode.GetEncodeDecoder(header.Encoding)
	err := decoder.Unmarshal(bodyBuf, val)
	if err != nil {
		return nil, log.NewErrorf("Service.callRouteFunc decoder.Unmarshal failed: %v", err)
	}
	// fmt.Printf("Service.callRouteFunc: %#v => %v\n", val, reflect.TypeOf(val))

	var result []interface{}
	aop.Recover(func() {
		result = method.Call(sess, helpers.GetValueFromPtr(val))
	}, func(e interface{}) {
		err = e.(error)
	})

	return result, err
}

func (s *ProcessorClient) response(sess *session.Session, header *pkg.Header, results []interface{}) error {
	respHeader := *header
	if header.Type == pkg.PKG_RPC_REQUEST {
		respHeader.Type = pkg.PKG_RPC_RESPONSE
	} else {
		respHeader.Type = pkg.PKG_RESPONSE
	}

	if results == nil || len(results) <= 0 {
		return sess.Send(&respHeader, []byte{})
	}

	result := results[0]
	/* check error != nil */
	if len(results) == 2 && !reflect.ValueOf(results[1]).IsNil() {
		// respHeader.Status = pkg.STAT_ERR
		respHeader.Status = result.(*pkg.ErrorMessage).Code
		if respHeader.Status == pkg.STAT_OK {
			log.Errorf("ErrorMessage.Code can't be STAT_OK!")
			respHeader.Status = pkg.STAT_ERR
		}
		result = results[1]
	}

	// fmt.Println("result:", result)

	encoder := encode.GetEncodeDecoder(header.Encoding)
	body, err := encoder.Marshal(result)
	if err != nil {
		return err
	}

	log.Logf("Send:\n\theader => %#v\n\tbody => %#v | %v", respHeader, body, string(body))
	return sess.Send(&respHeader, body)
}

func (s *ProcessorClient) recvPush(header *pkg.Header, body []byte) {
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
		item.Call(helpers.GetValueFromPtr(val))
	}
}

func (s *ProcessorClient) recvResponse(header *pkg.Header, body []byte) {
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
			cbs.successCallbak.Call(helpers.GetValueFromPtr(val))
			return
		}
	} else {
		val := cbs.failCallback.NewArg(0)
		err := s.Encoder.Unmarshal(body, val)
		if err == nil {
			cbs.failCallback.Call(helpers.GetValueFromPtr(val))
			return
		}
	}

	cbs.failCallback.Call(pkg.NewErrorMessage(
		pkg.STAT_ERR_DECODE_FAILED,
		fmt.Sprintf("decode body failed: %#v | %v", body, string(body))))
}

func (s *ProcessorClient) AddListener(route string, callback interface{}) {
	s.pushCbsMutex.Lock()
	defer s.pushCbsMutex.Unlock()

	list, ok := s.pushCbs[route]
	if !ok {
		list = make([]*Method, 0)
		s.pushCbs[route] = list
	}

	s.pushCbs[route] = append(s.pushCbs[route], NewMethod(callback))
}
