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
	"testing"

	"time"

	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
	"github.com/stretchr/testify/assert"
)

var callIn map[string]bool

type Message struct {
	Id  int
	Ok  bool
	M   map[string]int
	Arr []string
}

type ServerHandler struct {
	t *testing.T
}

func (self *ServerHandler) OnStarted() {
	callIn["OnStarted"] = true
	self.t.Logf("OnStarted %p\n", self)
}
func (self *ServerHandler) OnError(err error) {
	callIn["OnError"] = true
	self.t.Log("OnError", err)
}
func (self *ServerHandler) OnBeforeStop() {
	callIn["OnBeforeStop"] = true
	self.t.Log("OnBeforeStop")
}
func (self *ServerHandler) OnStopped() {
	callIn["OnStopped"] = true
	self.t.Log("OnStopped")
}
func (self *ServerHandler) OnNewClient(client transfer.IClient) {
	callIn["OnNewClient"] = true
	self.t.Log("OnNewClient", client)
	for {
		header, bodyBuf, err := client.Recv()
		callIn["Recv"] = true
		self.t.Log("Recv Error: ", err, header, bodyBuf, string(bodyBuf), len(bodyBuf))
		var obj Message
		err = encode.GetEncodeDecoder(pkg.ENCODING_JSON).Unmarshal(bodyBuf, &obj)
		self.t.Log("Decode Error: ", err)
		self.t.Logf("Recv:\nheader => %#v\nbodyBuf => %v | %v | %v \nmessage => %#v\n", header, bodyBuf, len(bodyBuf), string(bodyBuf), obj)
		if err != nil {
			break
		}
	}
}

func TestWebsocket(t *testing.T) {
	callIn = make(map[string]bool)

	serv := NewServer("localhost", 8888)
	serv.RegistDelegate(&ServerHandler{t})
	// go func() {
	err := serv.Start()
	assert.Nil(t, err)
	// }()

	cli := NewClient()
	err = cli.Connect("localhost", 8888)
	assert.Nil(t, err)

	header := cli.NewHeader(
		pkg.PKG_HEARTBEAT,
		pkg.ENCODING_JSON,
		"test.hello.world",
	)
	t.Log(header)
	obj := Message{
		Id: 1,
		Ok: true,
		M: map[string]int{
			"hello": 0,
			"world": 1,
		},
		Arr: []string{
			"from",
			"client",
		},
	}
	buf, err := encode.GetEncodeDecoder(pkg.ENCODING_JSON).Marshal(obj)
	assert.Nil(t, err, "Encode Error: %v", err)
	err = cli.Send(header, buf)
	assert.Nil(t, err)

	//Test Big Data
	buf = make([]byte, 4095)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 255)
	}
	err = cli.Send(header, buf)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)
	serv.Stop()

	assert.Equal(t, map[string]bool{
		"OnStarted":    true,
		"OnBeforeStop": true,
		"OnStopped":    true,
		"OnNewClient":  true,
		"Recv":         true,
	}, callIn)
}
