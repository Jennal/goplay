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
	"net"

	"github.com/jennal/goplay/transfer"
	"github.com/jennal/goplay/transfer/base"
)

type client struct {
	*base.Client
}

func NewClientWithConnect(conn net.Conn) transfer.IClient {
	ins := &client{
		Client: base.NewClientWithConnect(conn),
	}
	ins.SetImplement(ins)
	return ins
}

func NewClient() transfer.IClient {
	ins := &client{
		Client: base.NewClient(),
	}
	ins.SetImplement(ins)
	return ins
}

func (cli *client) CreateConn(host string, port int) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
}
