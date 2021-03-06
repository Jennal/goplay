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

//Package transfer defines how server and client connect to each other
package transfer

import (
	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/pkg"
)

type IHeaderCreator interface {
	NewHeader(t pkg.PackageType, e pkg.EncodingType, r string) *pkg.Header
	NewHeartBeatHeader() *pkg.Header
	NewHeartBeatResponseHeader(h *pkg.Header) *pkg.Header
}

type IClient interface {
	event.IEvent
	IHeaderCreator

	RegistDelegate(delegate IClientDelegate)
	UnregistDelegate(delegate IClientDelegate)

	LocalAddr() Addr
	RemoteAddr() Addr
	Id() uint32

	IsConnected() bool
	Connect(host string, port int) error
	Disconnect() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)

	Send(*pkg.Header, []byte) error
	Recv() (*pkg.Header, []byte, error)
}

type IClientDelegate interface {
	OnConnected(IClient)
	OnDisconnected(IClient)

	OnSent(IClient, *pkg.Header, []byte)
	OnRecved(IClient, *pkg.Header, []byte)
}

type IServer interface {
	event.IEvent

	RegistDelegate(delegate IServerDelegate)
	UnregistDelegate(delegate IServerDelegate)

	Host() string
	Port() int
	IsStarted() bool

	Addr() Addr
	Clients() map[uint32]IClient
	GetClientById(uint32) IClient

	Start() error
	Stop() error
}

type IServerDelegate interface {
	OnStarted()
	OnError(err error)
	OnBeforeStop()
	OnStopped()

	OnNewClient(client IClient)
}
