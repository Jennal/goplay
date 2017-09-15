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

package base

import (
	"net"

	"sync"

	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/transfer"
)

type IServerImplement interface {
	Open() error
	Accept() (transfer.IClient, error)
	Close() error

	Addr() net.Addr
}

type Server struct {
	*event.Event

	host      string
	port      int
	isStarted bool

	clients      map[uint32]transfer.IClient
	clientsMutex sync.Mutex

	exitChan chan bool
	impl     IServerImplement
}

func NewServer(host string, port int) *Server {
	return &Server{
		Event:     event.NewEvent(),
		host:      host,
		port:      port,
		clients:   make(map[uint32]transfer.IClient),
		isStarted: false,
		exitChan:  make(chan bool),
	}
}

func (server *Server) SetImplement(impl IServerImplement) {
	server.impl = impl
}

func (serv *Server) RegistDelegate(delegate transfer.IServerDelegate) {
	serv.On(transfer.EVENT_SERVER_STARTED, delegate, delegate.OnStarted)
	serv.On(transfer.EVENT_SERVER_BEFORE_STOP, delegate, delegate.OnBeforeStop)
	serv.On(transfer.EVENT_SERVER_STOPPED, delegate, delegate.OnStopped)
	serv.On(transfer.EVENT_SERVER_ERROR, delegate, delegate.OnError)
	serv.On(transfer.EVENT_SERVER_NEW_CLIENT, delegate, delegate.OnNewClient)
}

func (serv *Server) UnregistDelegate(delegate transfer.IServerDelegate) {
	serv.Off(transfer.EVENT_SERVER_STARTED, delegate)
	serv.Off(transfer.EVENT_SERVER_BEFORE_STOP, delegate)
	serv.Off(transfer.EVENT_SERVER_STOPPED, delegate)
	serv.Off(transfer.EVENT_SERVER_ERROR, delegate)
	serv.Off(transfer.EVENT_SERVER_NEW_CLIENT, delegate)
}

func (serv *Server) Host() string {
	return serv.host
}

func (serv *Server) Port() int {
	return serv.port
}

func (serv *Server) IsStarted() bool {
	return serv.isStarted
}

func (serv *Server) Addr() net.Addr {
	if serv.impl.Addr() == nil {
		return nil
	}

	return serv.impl.Addr()
}

func (serv *Server) Clients() map[uint32]transfer.IClient {
	return serv.clients
}

func (serv *Server) GetClientById(id uint32) transfer.IClient {
	serv.clientsMutex.Lock()
	defer serv.clientsMutex.Unlock()

	return serv.clients[id]
}

func (serv *Server) Start() error {
	if serv.isStarted {
		return nil
	}

	err := serv.impl.Open()
	if err != nil {
		return err
	}

	serv.isStarted = true

	go func() {
	Loop:
		for {
			select {
			case <-serv.exitChan:
				return
			default:

				client, err := serv.impl.Accept()
				if err != nil {
					if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
						continue Loop
					}

					serv.Emit(transfer.EVENT_SERVER_ERROR, err)
					serv.Stop()
					return
				}

				// fmt.Println("New Client:", conn.LocalAddr(), conn.RemoteAddr())
				client.Once(transfer.EVENT_CLIENT_DISCONNECTED, serv, func(c transfer.IClient) {
					serv.clientsMutex.Lock()
					defer serv.clientsMutex.Unlock()

					//remove from serv.clients
					for _, cli := range serv.clients {
						if c == cli {
							delete(serv.clients, cli.Id())
							return
						}
					}
				})
				serv.clientsMutex.Lock()
				serv.clients[client.Id()] = client
				serv.clientsMutex.Unlock()

				serv.Emit(transfer.EVENT_SERVER_NEW_CLIENT, client)
			}
		}
	}()

	defer serv.Emit(transfer.EVENT_SERVER_STARTED)
	return nil
}

func (serv *Server) Stop() error {
	if !serv.isStarted {
		return nil
	}

	serv.Emit(transfer.EVENT_SERVER_BEFORE_STOP)
	defer serv.Emit(transfer.EVENT_SERVER_STOPPED)
	serv.isStarted = false
	serv.exitChan <- true
	return serv.impl.Close()
}
