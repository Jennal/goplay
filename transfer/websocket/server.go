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

//Package websocket is a websocket implement to transfer
package websocket

import (
	"fmt"
	"net/http"

	"time"

	"github.com/gorilla/websocket"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/transfer"
	"github.com/jennal/goplay/transfer/base"
	"github.com/jennal/goplay/transfer/common"
)

const (
	URL_PREFIX = "ws://"
	URL_SUFFIX = "/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//Allow from all origin
		return true
	},
}

type server struct {
	*base.Server
	clientChan chan transfer.IClient
}

func NewServer(host string, port int) transfer.IServer {
	serv := &server{
		Server:     base.NewServer(host, port),
		clientChan: make(chan transfer.IClient),
	}
	serv.SetImplement(serv)
	return serv
}

func (serv *server) Open() error {
	log.Log("Open")
	host := fmt.Sprintf("%s:%d", serv.Host(), serv.Port())

	http.HandleFunc(URL_SUFFIX, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err)
			return
		}

		client := NewClientWithConnect(NewConn(conn))
		serv.clientChan <- client
	})

	log.Log("Open-1")
	var err error
	go func() {
		err = http.ListenAndServe(host, nil)
	}()
	log.Log("Open-2")
	log.Log("ListenAndServe:", err)
	if err != nil {
		return err
	}

	return nil
}

func (serv *server) Accept() (transfer.IClient, error) {
	select {
	case cli := <-serv.clientChan:
		log.Log("Accept")
		return cli, nil
	default:
		time.Sleep(100 * time.Millisecond)
		return nil, common.NewErrorTimeout("accept client timeout")
	}
}

func (serv *server) Close() error {
	/* DO NOTHING */
	return nil
}

func (serv *server) Addr() transfer.Addr {
	//TODO:
	return nil
}
