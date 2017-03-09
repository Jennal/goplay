package heartbeat

import (
	"time"

	"sync"

	"fmt"

	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer"
)

const (
	TIMEOUT     time.Duration = 3 * time.Second
	INTERNAL    time.Duration = 15 * time.Second
	MAX_TIMEOUT int           = 3
)

type HeartBeatFilter struct {
	times map[pkg.PackageIDType]time.Time
	mutex sync.Mutex

	lastPing  time.Duration
	avgPing   time.Duration
	minPing   time.Duration
	maxPing   time.Duration
	sendCount int64
	recvCount int64
}

func NewHeartBeatFilter() *HeartBeatFilter {
	result := &HeartBeatFilter{
		times: make(map[pkg.PackageIDType]time.Time),
	}
	go result.checkTimeOut()

	return result
}

func (self *HeartBeatFilter) pushTime(id pkg.PackageIDType, t time.Time) {
	self.mutex.Lock()
	self.mutex.Unlock()

	self.times[id] = t
}

func (self *HeartBeatFilter) popTime(id pkg.PackageIDType) *time.Time {
	self.mutex.Lock()
	self.mutex.Unlock()

	item, ok := self.times[id]
	if !ok {
		return nil
	}
	delete(self.times, id)

	return &item
}

func (self *HeartBeatFilter) checkTimeOut() {
	for {
		ids := []pkg.PackageIDType{}
		now := time.Now()

		for id, val := range self.times {
			if now.Sub(val) > TIMEOUT {
				ids = append(ids, id)
			}
		}

		for _, id := range ids {
			self.popTime(id)
		}

		time.Sleep(TIMEOUT)
	}
}

func (self *HeartBeatFilter) New(sess *session.Session) *pkg.Header {
	h := sess.NewHeartBeatHeader()
	self.pushTime(h.ID, time.Now())
	return h
}

func (self *HeartBeatFilter) OnNewClient(sess *session.Session) bool /* return false to ignore */ {
	exitSign := make(chan int)

	sess.On(transfer.EVENT_CLIENT_DISCONNECTED, self, func(cli transfer.IClient) {
		exitSign <- 1
	})

	go func() {
		for {
			select {
			case <-exitSign:
				return
			default:
				sess.Send(self.New(sess), []byte{})
				self.sendCount++
				time.Sleep(INTERNAL)
			}
		}
	}()

	return true
}

func (self *HeartBeatFilter) OnRecv(sess *session.Session, header *pkg.Header, body []byte) bool /* return false to ignore */ {
	if header.Type != pkg.PKG_HEARTBEAT && header.Type != pkg.PKG_HEARTBEAT_RESPONSE {
		return true
	}

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
		self.lastPing = time.Since(*lastTime)
		self.avgPing = (time.Duration)((float64(self.avgPing)*float64(self.recvCount-1) + float64(self.lastPing)) / float64(self.recvCount))
		if self.minPing == 0 || self.lastPing < self.minPing {
			self.minPing = self.lastPing
		}

		if self.maxPing == 0 || self.lastPing > self.maxPing {
			self.maxPing = self.lastPing
		}

		fmt.Println(self.Info())
	}

	return false
}

func (self *HeartBeatFilter) Info() string {
	return "Heart Beat Statistics ==========\n" +
		fmt.Sprintf("=> Last Ping: %v\n", self.lastPing/time.Millisecond) +
		fmt.Sprintf("=> Min Ping: %v\n", self.minPing/time.Millisecond) +
		fmt.Sprintf("=> Max Ping: %v\n", self.maxPing/time.Millisecond) +
		fmt.Sprintf("=> Avg Ping: %v\n", self.avgPing/time.Millisecond) +
		fmt.Sprintf("=> Send Count: %v\n", self.sendCount) +
		fmt.Sprintf("=> Recv Count: %v\n", self.recvCount) +
		fmt.Sprintf("=> Lost Rate: %.2g %%", 100-100*float64(self.recvCount)/float64(self.sendCount))
}
