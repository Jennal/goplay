package heartbeat

import (
	"fmt"
	"sync"
	"time"

	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

type HeartBeatManager struct {
	sync.Mutex
	items map[*session.Session]*HeartBeatProcessor

	lastPing  time.Duration
	avgPing   time.Duration
	minPing   time.Duration
	maxPing   time.Duration
	sendCount int64
	recvCount int64
}

func NewHeartBeatManager() *HeartBeatManager {
	return &HeartBeatManager{
		items: make(map[*session.Session]*HeartBeatProcessor),
	}
}

func (self *HeartBeatManager) OnNewClient(sess *session.Session) bool /* return false to ignore */ {
	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		self.Lock()
		defer self.Unlock()

		if _, ok := self.items[sess]; ok {
			delete(self.items, sess)
		}
	})

	self.Lock()
	defer self.Unlock()

	f := NewHeartBeatProcessor(self)
	f.OnNewClient(sess)
	self.items[sess] = f

	return true
}

func (self *HeartBeatManager) OnRecv(sess *session.Session, header *pkg.Header, body []byte) bool /* return false to ignore */ {
	if header.Type != pkg.PKG_HEARTBEAT && header.Type != pkg.PKG_HEARTBEAT_RESPONSE {
		return true
	}

	self.Lock()
	f, ok := self.items[sess]
	self.Unlock()
	if !ok {
		return true
	}

	result := f.OnRecv(sess, header, body)
	self.summary(f)
	// fmt.Println(self.Info())
	return result
}

func (self *HeartBeatManager) summary(f *HeartBeatProcessor) {
	self.lastPing = f.lastPing
	self.avgPing = (time.Duration)((float64(self.avgPing)*float64(self.recvCount-1) + float64(self.lastPing)) / float64(self.recvCount))
	if self.minPing == 0 || self.lastPing < self.minPing {
		self.minPing = self.lastPing
	}

	if self.maxPing == 0 || self.lastPing > self.maxPing {
		self.maxPing = self.lastPing
	}
}

func (self *HeartBeatManager) incSendCount() {
	self.sendCount++
}

func (self *HeartBeatManager) incRecvCount() {
	self.recvCount++
}

func (self *HeartBeatManager) Info() string {
	return "Heart Beat Summary Statistics ==========\n" +
		fmt.Sprintf("=> Last Ping: %v\n", self.lastPing/time.Millisecond) +
		fmt.Sprintf("=> Min Ping: %v\n", self.minPing/time.Millisecond) +
		fmt.Sprintf("=> Max Ping: %v\n", self.maxPing/time.Millisecond) +
		fmt.Sprintf("=> Avg Ping: %v\n", self.avgPing/time.Millisecond) +
		fmt.Sprintf("=> Send Count: %v\n", self.sendCount) +
		fmt.Sprintf("=> Recv Count: %v\n", self.recvCount) +
		fmt.Sprintf("=> Lost Rate: %.2g %%", 100-100*float64(self.recvCount)/float64(self.sendCount))
}
