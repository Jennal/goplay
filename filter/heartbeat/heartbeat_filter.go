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

package heartbeat

import (
	"fmt"
	"sync"
	"time"

	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

const (
	TIMEOUT     time.Duration = 3 * time.Second
	INTERNAL    time.Duration = 15 * time.Second
	MAX_TIMEOUT int           = 3
)

type HeartBeatProcessor struct {
	mapMutex     sync.Mutex
	timeoutMutex sync.Mutex

	manager *HeartBeatManager
	sess    *session.Session
	times   map[pkg.PackageIDType]time.Time

	lastPing     time.Duration
	avgPing      time.Duration
	minPing      time.Duration
	maxPing      time.Duration
	sendCount    int64
	recvCount    int64
	timeoutCount int
}

func NewHeartBeatProcessor(manager *HeartBeatManager) *HeartBeatProcessor {
	result := &HeartBeatProcessor{
		manager: manager,
		times:   make(map[pkg.PackageIDType]time.Time),
	}
	go result.checkTimeOut()

	return result
}

func (self *HeartBeatProcessor) pushTime(id pkg.PackageIDType, t time.Time) {
	self.mapMutex.Lock()
	defer self.mapMutex.Unlock()

	self.times[id] = t
}

func (self *HeartBeatProcessor) popTime(id pkg.PackageIDType) *time.Time {
	self.mapMutex.Lock()
	defer self.mapMutex.Unlock()

	item, ok := self.times[id]
	if !ok {
		return nil
	}
	delete(self.times, id)

	return &item
}

func (self *HeartBeatProcessor) checkTimeOut() {
	//FIXME: this will lead to memory leak
	for {
		// fmt.Printf("[%v] Check Timeout\n", time.Now())
		ids := []pkg.PackageIDType{}
		now := time.Now()

		for id, val := range self.times {
			if now.Sub(val) > TIMEOUT {
				// fmt.Printf("\t => Timeout: %v => %v, id = %v\n", now, val, id)
				ids = append(ids, id)
			}
		}

		for _, id := range ids {
			self.popTime(id)
			if !self.incTimeOut() {
				return
			}
		}

		time.Sleep(TIMEOUT)
	}
}

func (self *HeartBeatProcessor) incTimeOut() bool {
	self.timeoutMutex.Lock()
	defer self.timeoutMutex.Unlock()

	self.timeoutCount++
	// fmt.Println("TimeOut:", self.sess, self.timeoutCount, MAX_TIMEOUT)
	if self.timeoutCount >= MAX_TIMEOUT {
		self.sess.Disconnect()
		return false
	}

	return true
}

func (self *HeartBeatProcessor) resetTimeOut() {
	self.timeoutMutex.Lock()
	defer self.timeoutMutex.Unlock()

	self.timeoutCount = 0
	// fmt.Println("Reset timeOut:", self.sess, self.timeoutCount, MAX_TIMEOUT)
}

func (self *HeartBeatProcessor) New(sess *session.Session) *pkg.Header {
	h := sess.NewHeartBeatHeader()
	self.pushTime(h.ID, time.Now())
	return h
}

func (self *HeartBeatProcessor) OnNewClient(sess *session.Session) bool /* return false to ignore */ {
	self.sess = sess
	exitSign := make(chan int)

	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		// fmt.Println("HeartBeatProcessor-Client-Disconnected", sess)
		exitSign <- 1
	})

	go func() {
		for {
			select {
			case <-exitSign:
				return
			default:
				if err := sess.Send(self.New(sess), []byte{}); err == nil {
					self.sendCount++
					self.manager.incSendCount()
				} else {
					self.incTimeOut()
				}

				time.Sleep(INTERNAL)
			}
		}
	}()

	return true
}

func (self *HeartBeatProcessor) OnRecv(sess *session.Session, header *pkg.Header, body []byte) bool /* return false to ignore */ {
	if header.Type != pkg.PKG_HEARTBEAT && header.Type != pkg.PKG_HEARTBEAT_RESPONSE {
		return true
	}

	if sess != self.sess {
		log.Errorf("can't come to here, must be something wrong")
		return true
	}

	self.resetTimeOut()

	switch header.Type {
	case pkg.PKG_HEARTBEAT:
		resp := sess.NewHeartBeatResponseHeader(header)
		sess.Send(resp, []byte{})
	case pkg.PKG_HEARTBEAT_RESPONSE:
		lastTime := self.popTime(header.ID)
		if lastTime == nil {
			return false
		}

		self.recvCount++
		self.manager.incRecvCount()
		self.lastPing = time.Since(*lastTime)
		self.avgPing = (time.Duration)((float64(self.avgPing)*float64(self.recvCount-1) + float64(self.lastPing)) / float64(self.recvCount))
		if self.minPing == 0 || self.lastPing < self.minPing {
			self.minPing = self.lastPing
		}

		if self.maxPing == 0 || self.lastPing > self.maxPing {
			self.maxPing = self.lastPing
		}

		// fmt.Println(self.Info())
	}

	return false
}

func (self *HeartBeatProcessor) Info() string {
	return fmt.Sprintf("Heart Beat For Sess(%d) Statistics ==========\n", self.sess.ID) +
		fmt.Sprintf("=> Last Ping: %v\n", self.lastPing/time.Millisecond) +
		fmt.Sprintf("=> Min Ping: %v\n", self.minPing/time.Millisecond) +
		fmt.Sprintf("=> Max Ping: %v\n", self.maxPing/time.Millisecond) +
		fmt.Sprintf("=> Avg Ping: %v\n", self.avgPing/time.Millisecond) +
		fmt.Sprintf("=> Send Count: %v\n", self.sendCount) +
		fmt.Sprintf("=> Recv Count: %v\n", self.recvCount) +
		fmt.Sprintf("=> Lost Rate: %.2g %%", 100-100*float64(self.recvCount)/float64(self.sendCount))
}
