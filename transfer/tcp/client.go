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

package tcp

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
	"github.com/jennal/goplay/transfer/base"
)

var (
	ERR_ALREADY_CONNECTED      = errors.New("already connected")
	ERR_READ_BEFORE_CONNECTED  = errors.New("read before connected")
	ERR_WRITE_BEFORE_CONNECTED = errors.New("write before connected")
)

var idGen = helpers.NewIDGen(defaults.MAX_CLIENT_COUNT)

type client struct {
	*event.Event
	*base.HeaderCreator

	conn        net.Conn
	isConnected bool
	id          uint32

	sendMutex sync.Mutex
	recvMutex sync.Mutex
}

func NewClientWithConnect(conn net.Conn) transfer.IClient {
	return &client{
		Event:         event.NewEvent(),
		HeaderCreator: base.NewHeaderCreator(),
		conn:          conn,
		isConnected:   true,
		id:            idGen.NextID(),
	}
}

func NewClient() transfer.IClient {
	return &client{
		Event:         event.NewEvent(),
		HeaderCreator: base.NewHeaderCreator(),
		isConnected:   false,
		id:            idGen.NextID(),
	}
}

func (client *client) RegistDelegate(delegate transfer.IClientDelegate) {
	client.On(transfer.EVENT_CLIENT_CONNECTED, delegate, delegate.OnConnected)
	client.On(transfer.EVENT_CLIENT_DISCONNECTED, delegate, delegate.OnDisconnected)
	client.On(transfer.EVENT_CLIENT_SENT, delegate, delegate.OnSent)
	client.On(transfer.EVENT_CLIENT_RECVED, delegate, delegate.OnRecved)
}

func (client *client) UnregistDelegate(delegate transfer.IClientDelegate) {
	client.Off(transfer.EVENT_CLIENT_CONNECTED, delegate)
	client.Off(transfer.EVENT_CLIENT_DISCONNECTED, delegate)
	client.Off(transfer.EVENT_CLIENT_SENT, delegate)
	client.Off(transfer.EVENT_CLIENT_RECVED, delegate)
}

func (client *client) LocalAddr() net.Addr {
	if client.conn == nil {
		return nil
	}

	return client.conn.LocalAddr()
}

func (client *client) RemoteAddr() net.Addr {
	if client.conn == nil {
		return nil
	}

	return client.conn.RemoteAddr()
}

func (client *client) Id() uint32 {
	return client.id
}

func (client *client) IsConnected() bool {
	return client.conn != nil && client.isConnected
}

func (client *client) Connect(host string, port int) error {
	if client.isConnected {
		return ERR_ALREADY_CONNECTED
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	client.conn = conn
	client.isConnected = true

	defer client.Emit(transfer.EVENT_CLIENT_CONNECTED, client)
	return nil
}

func (client *client) Disconnect() error {
	if !client.IsConnected() {
		return nil
	}

	defer client.Emit(transfer.EVENT_CLIENT_DISCONNECTED, client)
	// log.Logf("************ Disconnectd: %p => %#v", client, client.Event)
	client.isConnected = false
	return client.conn.Close()
}

func (client *client) Read(buf []byte) (int, error) {
	if !client.IsConnected() {
		return 0, ERR_READ_BEFORE_CONNECTED
	}

	size := len(buf)
	n := 0

	for n < size {
		rn, err := client.conn.Read(buf[n:])
		n += rn
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (client *client) Write(buf []byte) (int, error) {
	if !client.IsConnected() {
		return 0, ERR_WRITE_BEFORE_CONNECTED
	}

	size := len(buf)
	n := 0
	for n < size {
		wn, err := client.conn.Write(buf[n:])
		n += wn
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (client *client) Send(header *pkg.Header, data []byte) error {
	client.sendMutex.Lock()
	defer client.sendMutex.Unlock()

	header.ContentSize = pkg.PackageSizeType(len(data))
	headerBuffer, err := header.Marshal()
	if err != nil {
		return err
	}
	buffer := append(headerBuffer, data...)
	// log.Logf("Write:\n\theader => %#v\n\tbody => %#v | %v\n", header, data, string(data))

	_, err = client.Write(buffer)

	if err == nil {
		defer client.Emit(transfer.EVENT_CLIENT_SENT, client, header, data)
	}

	return err
}

func (client *client) Recv() (*pkg.Header, []byte, error) {
	client.recvMutex.Lock()
	defer client.recvMutex.Unlock()

	header := &pkg.Header{}
	_, err := pkg.ReadHeader(client, header)
	if err != nil {
		return nil, nil, err
	}

	if header.ContentSize > 0 {
		buffer := make([]byte, header.ContentSize)
		_, err := client.Read(buffer)
		if err != nil {
			return nil, nil, err
		}

		defer client.Emit(transfer.EVENT_CLIENT_RECVED, client, header, buffer)
		// fmt.Println("Recv body:", buffer)
		return header, buffer, err
	}

	defer client.Emit(transfer.EVENT_CLIENT_RECVED, client, header, nil)
	return header, nil, nil
}
