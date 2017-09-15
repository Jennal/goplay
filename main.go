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

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jennal/goplay/aop"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
	"github.com/jennal/goplay/transfer/websocket"
)

func init() {
	fmt.Println("init-1")
}

func init() {
	fmt.Println("init-2")
}

type Handler struct{}

func (self *Handler) OnStarted() {
	fmt.Println("Handler-OnStarted")
}

func (self *Handler) OnError(err error) {
	fmt.Println("OnError", err)
}

func (self *Handler) OnBeforeStop() {
	fmt.Println("OnBeforeStop")
}

func (self *Handler) OnStopped() {
	fmt.Println("OnStopped")
}

func (self *Handler) OnNewClient(cli transfer.IClient) {
	fmt.Println("Handler-OnNewClient", cli)

	go func() {
		aop.Recover(func() {
			for {
				header, body, err := cli.Recv()
				log.Logf("Recv:\n\theader => %#v\n\tbody => %#v | %v\n\terr => %v\n", header, body, string(body), err)
				err = cli.Send(header, body)
				if err != nil {
					log.Error(err)
				}
			}
		}, func(err interface{}) {
			log.RecoverErrorf("error: %#v", err)
		})
	}()
}

func (self *Handler) Test(sess *session.Session, line string) *pkg.ErrorMessage {
	fmt.Println("Handler-Test", sess, line)
	sess.Push("test.push", "Service: "+line)
	return nil
}

func (self *Handler) Add(sess *session.Session, n int) (int, *pkg.ErrorMessage) {
	fmt.Println("Handler-Add", sess, n)
	return n + 1, nil
}

func main() {
	str := ""

	ser := websocket.NewServer("localhost", 9990)
	ser.RegistDelegate(&Handler{})
	// serv := service.NewService("test", ser)

	// serv.RegistHanlder(&Handler{})

	err := ser.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("server started:", ser.Addr())
	fmt.Scanf("%s", &str)

	cli := websocket.NewClient()
	cli.Connect("localhost", 9990)

	idGen := pkg.NewIDGen()
	scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println(scanner)
	for scanner.Scan() {
		// fmt.Println(scanner)
		buffer := scanner.Text()
		fmt.Println("Send: ", buffer)
		cli.Send(pkg.NewHeader(pkg.PKG_NOTIFY, pkg.ENCODING_JSON, idGen, "test.route"), []byte(buffer))
	}

	// client := service.NewServiceClient(cli)
	// client.AddListener("test.push", func(line string) {
	// 	fmt.Println("[test.push] recv: ", line)
	// })
	// client.Connect("", 9990)

	// client.Request("test.handler.add", 1, func(result int) {
	// 	fmt.Println("[test.handler.add] callback: ", result)
	// }, func(err *pkg.ErrorMessage) {
	// 	fmt.Println("[test.handler.add] error: ", err)
	// })

	// client.Notify("test.handler.test", "Hello from Client")

	// data, err := json.Marshal("Hello From Client")
	// fmt.Println("json encode:", string(data))
	// if err != nil {
	// 	fmt.Println("json.Marshal error:", err)
	// } else {
	// 	// var str string
	// 	// json.Unmarshal(data, &str)
	// 	// fmt.Println("Unmarshal:", str)
	// 	cli.Send(cli.NewHeader(pkg.PKG_NOTIFY, pkg.ENCODING_JSON, "test.handler.test"), data)
	// }

	// data, err = json.Marshal(1000)
	// fmt.Println("json encode:", string(data))
	// if err != nil {
	// 	fmt.Println("json.Marshal error:", err)
	// } else {
	// 	cli.Send(cli.NewHeader(pkg.PKG_REQUEST, pkg.ENCODING_JSON, "test.handler.add"), data)
	// }

	// go func() {
	// 	for i := 0; true; i++ {
	// 		header, body, err := cli.Recv()
	// 		fmt.Printf("[%v] Recv:\n\theader => %#v\n\tbody => %#v | %v\n\terr => %v\n", time.Now(), header, body, string(body), err)
	// 		if err != nil {
	// 			break
	// 		}
	// 		if /*i%4 == 0 &&*/ header.Type == pkg.PKG_HEARTBEAT {
	// 			respHeader := cli.NewHeartBeatResponseHeader(header)
	// 			cli.Send(respHeader, []byte{})
	// 		}
	// 	}
	// }()

	// fmt.Scanf("%s", nil)
	// ser.Stop()
	// fmt.Scanf("%s", nil)
}
