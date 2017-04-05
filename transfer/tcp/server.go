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

//Package tcp is a tcp implement to transfer
package tcp

import (
	"fmt"
	"net"

	"sync"

	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/transfer"
)

type server struct {
	*event.Event

	host      string
	port      int
	isStarted bool

	clients      map[uint32]transfer.IClient
	clientsMutex sync.Mutex

	listener net.Listener
}

func NewServer(host string, port int) transfer.IServer {
	return &server{
		Event:     event.NewEvent(),
		host:      host,
		port:      port,
		clients:   make(map[uint32]transfer.IClient),
		isStarted: false,
	}
}

func (serv *server) RegistDelegate(delegate transfer.IServerDelegate) {
	serv.On(transfer.EVENT_SERVER_STARTED, delegate, delegate.OnStarted)
	serv.On(transfer.EVENT_SERVER_BEFORE_STOP, delegate, delegate.OnBeforeStop)
	serv.On(transfer.EVENT_SERVER_STOPPED, delegate, delegate.OnStopped)
	serv.On(transfer.EVENT_SERVER_ERROR, delegate, delegate.OnError)
	serv.On(transfer.EVENT_SERVER_NEW_CLIENT, delegate, delegate.OnNewClient)
}

func (serv *server) UnregistDelegate(delegate transfer.IServerDelegate) {
	serv.Off(transfer.EVENT_SERVER_STARTED, delegate)
	serv.Off(transfer.EVENT_SERVER_BEFORE_STOP, delegate)
	serv.Off(transfer.EVENT_SERVER_STOPPED, delegate)
	serv.Off(transfer.EVENT_SERVER_ERROR, delegate)
	serv.Off(transfer.EVENT_SERVER_NEW_CLIENT, delegate)
}

func (serv *server) Host() string {
	return serv.host
}

func (serv *server) Port() int {
	return serv.port
}

func (serv *server) IsStarted() bool {
	return serv.isStarted
}

func (serv *server) Addr() net.Addr {
	if serv.listener == nil {
		return nil
	}

	return serv.listener.Addr()
}

func (serv *server) Clients() map[uint32]transfer.IClient {
	return serv.clients
}

func (serv *server) GetClientById(id uint32) transfer.IClient {
	serv.clientsMutex.Lock()
	defer serv.clientsMutex.Unlock()

	return serv.clients[id]
}

func (serv *server) Start() error {
	if serv.isStarted {
		return nil
	}

	host := fmt.Sprintf("%s:%d", serv.host, serv.port)
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	serv.isStarted = true

	serv.listener = ln
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				serv.Emit(transfer.EVENT_SERVER_ERROR, err)
				serv.Stop()
				return
			}

			// fmt.Println("New Client:", conn.LocalAddr(), conn.RemoteAddr())
			client := NewClientWithConnect(conn)
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
	}()

	defer serv.Emit(transfer.EVENT_SERVER_STARTED)
	return nil
}

func (serv *server) Stop() error {
	if !serv.isStarted {
		return nil
	}

	serv.Emit(transfer.EVENT_SERVER_BEFORE_STOP)
	defer serv.Emit(transfer.EVENT_SERVER_STOPPED)
	serv.isStarted = false
	return serv.listener.Close()
}
