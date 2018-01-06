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

package channel

import (
	"fmt"
	"strings"

	"github.com/jennal/goplay/data"
	"github.com/jennal/goplay/session"
)

const (
	CHANNEL_PREFIX = "channel."
)

type Channel struct {
	*data.Map
	*session.SessionManager
	name string
}

func NewChannel(name string) *Channel {
	return &Channel{
		Map:            data.NewMap(),
		SessionManager: session.NewSessionManager(),
		name:           name,
	}
}

func (ch *Channel) Name() string {
	return ch.name
}

func (ch *Channel) Route(route string) string {
	route = strings.ToLower(route)
	if strings.HasPrefix(route, CHANNEL_PREFIX) {
		return route
	}

	return fmt.Sprintf("%v%v", CHANNEL_PREFIX, route)
}

//Broadcast direct to client
func (ch *Channel) Broadcast(route string, obj interface{}) {
	route = ch.Route(route)
	for _, sess := range ch.Sessions() {
		sess.Push(route, obj)
	}
}

//Broadcast across backend
func (ch *Channel) BroadcastRaw(route string, data []byte) {
	route = ch.Route(route)
	for _, sess := range ch.Sessions() {
		sess.Broadcast(route, data)
	}
}
