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
	"time"

	"github.com/jennal/goplay/transfer"
	"github.com/jennal/goplay/transfer/base"
)

type server struct {
	*base.Server
	listener *net.TCPListener
}

func NewServer(host string, port int) transfer.IServer {
	serv := &server{
		Server: base.NewServer(host, port),
	}
	serv.SetImplement(serv)
	return serv
}

func (serv *server) Open() error {
	host := fmt.Sprintf("%s:%d", serv.Host(), serv.Port())
	laddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	serv.listener = ln
	return nil
}

func (serv *server) Accept() (transfer.IClient, error) {
	serv.listener.SetDeadline(time.Now().Add(time.Second))
	conn, err := serv.listener.Accept()
	if err != nil {
		return nil, err
	}

	// fmt.Println("New Client:", conn.LocalAddr(), conn.RemoteAddr())
	client := NewClientWithConnect(conn)
	return client, nil
}

func (serv *server) Close() error {
	return serv.listener.Close()
}

func (serv *server) Addr() transfer.Addr {
	if serv.listener == nil {
		return nil
	}

	return transfer.NewAddr(serv.listener.Addr())
}
