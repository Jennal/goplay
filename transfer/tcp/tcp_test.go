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
	"fmt"
	"testing"

	"time"

	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
	"github.com/stretchr/testify/assert"
)

type Message struct {
	Id  int
	Ok  bool
	M   map[string]int
	Arr []string
}

type ServerHandler struct {
}

func (self *ServerHandler) OnStarted() {
	fmt.Printf("OnStarted %p\n", self)
}
func (self *ServerHandler) OnError(err error) {
	fmt.Println("OnError", err)
}
func (self *ServerHandler) OnBeforeStop() {
	fmt.Println("OnBeforeStop")
}
func (self *ServerHandler) OnStopped() {
	fmt.Println("OnStopped")
}
func (self *ServerHandler) OnNewClient(client transfer.IClient) {
	fmt.Println("OnNewClient", client)
	for {
		header, bodyBuf, err := client.Recv()
		fmt.Println("Recv Error: ", err, header, bodyBuf, string(bodyBuf))
		var obj Message
		err = encode.GetEncodeDecoder(pkg.ENCODING_JSON).Unmarshal(bodyBuf, &obj)
		fmt.Println("Decode Error: ", err)
		fmt.Printf("Recv:\nheader => %#v\nbodyBuf => %v\nmessage => %#v\n", header, bodyBuf, obj)
		if err != nil {
			break
		}
	}
}

func TestTcp(t *testing.T) {
	serv := NewServer("", 8888)
	serv.RegistDelegate(&ServerHandler{})
	go serv.Start()

	cli := NewClient()
	err := cli.Connect("", 8888)
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
	cli.Send(header, buf)
	time.Sleep(time.Second)
}
