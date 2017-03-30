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

package session

import (
	"sync"

	"github.com/jennal/goplay/transfer"
)

type SessionManager struct {
	sync.Mutex
	sessions []*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make([]*Session, 0),
	}
}

func (self *SessionManager) Add(sess *Session) {
	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(client transfer.IClient) {
		self.Remove(sess)
	})

	self.Lock()
	defer self.Unlock()
	self.sessions = append(self.sessions, sess)
}

func (self *SessionManager) Remove(sess *Session) {
	self.Lock()
	defer self.Unlock()

	for i, s := range self.sessions {
		if s == sess {
			self.sessions = append(self.sessions[:i], self.sessions[i+1:]...)
			return
		}
	}
}

func (self *SessionManager) Sessions() []*Session {
	return self.sessions
}
