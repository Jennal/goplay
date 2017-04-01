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
	sessions map[uint32]map[uint32]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[uint32]map[uint32]*Session),
	}
}

func (self *SessionManager) Add(sess *Session) {
	if self.Exists(sess) {
		/* already added */
		return
	}

	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(client transfer.IClient) {
		self.Remove(sess)
	})

	self.Lock()
	defer self.Unlock()

	if _, ok := self.sessions[sess.ID]; !ok {
		self.sessions[sess.ID] = make(map[uint32]*Session)
	}

	self.sessions[sess.ID][sess.ClientID] = sess
}

func (self *SessionManager) Remove(sess *Session) {
	self.Lock()
	defer self.Unlock()

	if subMap, ok := self.sessions[sess.ID]; ok {
		if _, ok := subMap[sess.ClientID]; ok {
			delete(self.sessions[sess.ID], sess.ClientID)
		}
	}
}

func (self *SessionManager) Exists(sess *Session) bool {
	s := self.GetSessionByID(sess.ID, sess.ClientID)
	return s != nil
}

func (self *SessionManager) Count() int {
	count := 0

	self.Lock()
	defer self.Unlock()

	for _, subMap := range self.sessions {
		count += len(subMap)
	}

	return count
}

func (self *SessionManager) Sessions() []*Session {
	result := make([]*Session, 0)

	self.Lock()
	defer self.Unlock()

	for _, subMap := range self.sessions {
		for _, sess := range subMap {
			result = append(result, sess)
		}
	}

	return result
}

func (self *SessionManager) GetSessionByID(id uint32, clientId uint32) *Session {
	self.Lock()
	defer self.Unlock()

	if subMap, ok := self.sessions[id]; ok {
		if sess, ok := subMap[clientId]; ok {
			return sess
		}
	}

	return nil
}
