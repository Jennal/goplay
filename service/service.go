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

//Package service is main loop logic
package service

import (
	"sync"

	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/filter/handshake"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type Service struct {
	transfer.IServer
	*SettingContainer
	router *router.Router

	Name     string
	Encoding pkg.EncodingType

	handlers []handler.IHandler
	filters  []filter.IFilter

	heartBeatManager *heartbeat.HeartBeatManager

	clients      []*ProcessorClient
	clientsMutex sync.Mutex
}

func NewService(name string, serv transfer.IServer) *Service {
	instance := &Service{
		Name:     name,
		Encoding: defaults.Encoding,

		SettingContainer: NewSettingContainer(),
		IServer:          serv,
		router:           router.NewRouter(),
		heartBeatManager: heartbeat.NewHeartBeatManager(),
	}

	serv.RegistDelegate(instance)
	instance.RegistFilter(handshake.NewHandShakeFilter())

	return instance
}

func (self *Service) Router() *router.Router {
	return self.router
}

func (self *Service) Handlers() []handler.IHandler {
	return self.handlers
}

func (self *Service) Filters() []filter.IFilter {
	return self.filters
}

func (self *Service) ServiceClients() []*ProcessorClient {
	return self.clients
}

func (self *Service) SetEncoding(e pkg.EncodingType) error {
	if encoder := encode.GetEncodeDecoder(e); encoder != nil {
		self.Encoding = e
		return nil
	}

	return log.NewErrorf("can't find encoder with: %v", e)
}

func (self *Service) SetSettings(s *Settings) error {
	err := self.SettingContainer.SetSettings(s)
	if err != nil {
		return err
	}

	self.clientsMutex.Lock()
	defer self.clientsMutex.Unlock()

	for _, c := range self.clients {
		c.SetSettings(s)
	}
	return nil
}

func (self *Service) RegistHanlder(obj handler.IHandler) {
	self.router.Add(self.Name, obj)
	self.handlers = append(self.handlers, obj)
}

func (self *Service) RegistHanlderGroup(group map[string][]handler.IHandler) {
	for name, list := range group {
		for _, h := range list {
			self.router.Add(name, h)
			self.handlers = append(self.handlers, h)
		}
	}
}

func (self *Service) RegistFilter(obj filter.IFilter) {
	self.filters = append(self.filters, obj)
}

func (self *Service) OnStarted() {
	log.Logf("OnStarted %p", self)
	for _, handler := range self.handlers {
		handler.OnStarted()
	}
}

func (self *Service) OnError(err error) {
	log.Error(err)
}

func (self *Service) OnBeforeStop() {
	log.Log("OnBeforeStop")
	self.clientsMutex.Lock()
	defer self.clientsMutex.Unlock()

	for _, client := range self.clients {
		client.Push(ON_SERVICE_DOWN, true)
	}
	// time.Sleep(100 * time.Millisecond)
}

func (self *Service) OnStopped() {
	log.Log("OnStopped")
	for _, handler := range self.handlers {
		handler.OnStopped()
	}
}

func (self *Service) RegistNewClient(client transfer.IClient) *ProcessorClient {
	log.Log("OnNewClient:", client)
	serviceClient := NewProcessorClient(client)
	serviceClient.SetSettings(self.Settings())
	serviceClient.SetEncoding(self.Encoding)
	serviceClient.SetRouter(self.router)
	serviceClient.SetFilters(self.filters)
	serviceClient.SetHeartBeatManager(self.heartBeatManager)

	serviceClient.Bind(session.IDGen.NextID())

	serviceClient.Once(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		self.clientsMutex.Lock()
		defer self.clientsMutex.Unlock()

		for i, client := range self.clients {
			if client == serviceClient {
				self.clients = append(self.clients[:i], self.clients[i+1:]...)
				break
			}
		}
	})
	self.clientsMutex.Lock()
	self.clients = append(self.clients, serviceClient)
	self.clientsMutex.Unlock()
	// log.Log(len(self.clients), self.clients)

	return serviceClient
}

func (self *Service) FilterOnNewClient(sess *session.Session) bool {
	for _, filter := range self.filters {
		if !filter.OnNewClient(sess) {
			return false
		}
	}

	return true
}

func (self *Service) HandlerOnNewClient(sess *session.Session) {
	for _, handler := range self.handlers {
		handler.OnNewClient(sess)
	}
}

func (self *Service) OnNewClient(client transfer.IClient) {
	serviceClient := self.RegistNewClient(client)
	sess := serviceClient.getSession(serviceClient.ID, serviceClient.ClientID)
	if self.FilterOnNewClient(sess) {
		self.HandlerOnNewClient(sess)
		serviceClient.Emit(transfer.EVENT_CLIENT_CONNECTED, client)
	}
}
