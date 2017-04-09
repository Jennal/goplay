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

package websocket

import (
	"errors"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ERR_WRONG_MESSAGE_TYPE = errors.New("wrong message type")
)

type Conn struct {
	*websocket.Conn
}

func NewConn(conn *websocket.Conn) net.Conn {
	return &Conn{
		Conn: conn,
	}
}

func (conn *Conn) Read(b []byte) (int, error) {
	t, buffer, err := conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	if t != websocket.BinaryMessage {
		return 0, ERR_WRONG_MESSAGE_TYPE
	}

	copy(b, buffer)
	return len(buffer), nil
}

func (conn *Conn) Write(b []byte) (int, error) {
	err := conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (conn *Conn) SetDeadline(t time.Time) error {
	err := conn.SetReadDeadline(t)
	if err != nil {
		return err
	}
	err = conn.SetWriteDeadline(t)
	if err != nil {
		return err
	}

	return nil
}
