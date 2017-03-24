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
	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/filter"
	"github.com/jennal/goplay/filter/heartbeat"
	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
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
	log.Logf("OnStarted %p", self)
	for _, handler := range self.handlers {
		handler.OnStarted()
	}
}

func (self *Service) OnError(err error) {
	log.Error(err)
}

func (self *Service) OnStopped() {
	log.Log("OnStopped")
	for _, handler := range self.handlers {
		handler.OnStopped()
	}
}

func (self *Service) OnNewClient(client transfer.IClient) {
	log.Log("OnNewClient:", client)
	serviceClient := NewServiceClient(client)
	serviceClient.SetEncoding(self.Encoding)
	serviceClient.SetRouter(self.router)
	serviceClient.SetFilters(self.filters)

	for _, handler := range self.handlers {
		handler.OnNewClient(serviceClient.Session)
	}

	serviceClient.Emit(transfer.EVENT_CLIENT_CONNECTED, client)
}
