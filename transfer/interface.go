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
	"net"

	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer/base"
)

type IClient interface {
	event.IEvent
	base.IHeaderCreator

	RegistDelegate(delegate IClientDelegate)
	UnregistDelegate(delegate IClientDelegate)

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

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
}

type IServer interface {
	event.IEvent

	RegistDelegate(delegate IServerDelegate)
	UnregistDelegate(delegate IServerDelegate)

	Addr() net.Addr
	Clients() []IClient

	Start() error
	Stop() error
}

type IServerDelegate interface {
	OnStarted()
	OnError(err error)
	OnStopped()

	OnNewClient(client IClient)
}
