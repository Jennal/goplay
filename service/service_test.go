package service

import (
	"testing"

	"time"

	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
	"github.com/jennal/goplay/transfer/tcp"
	"github.com/stretchr/testify/assert"
)

const (
	PORT     = 9000
	Encoding = defaults.Encoding
)

type Handler struct {
	t           *testing.T
	Started     bool
	Stopped     bool
	ClientCount int
}

func NewHandler(t *testing.T) *Handler {
	return &Handler{
		t:           t,
		Started:     false,
		Stopped:     false,
		ClientCount: 0,
	}
}

func (self *Handler) OnStarted() {
	self.t.Log("Handler-OnStarted")
	self.Started = true
}

func (self *Handler) OnStopped() {
	self.t.Log("Handler-OnStopped")
	self.Stopped = true
}

func (self *Handler) OnNewClient(sess *session.Session) {
	self.t.Log("Handler-OnNewClient", sess)
	self.ClientCount++
}

func (self *Handler) NotifyString(sess *session.Session, line string) *pkg.ErrorMessage {
	self.t.Log("Handler-Test", sess, line)
	assert.Equal(self.t, "string", line)
	sess.Push("test.pushstring", "Service: "+line)
	return nil
}

func (self *Handler) RequestInt(sess *session.Session, n int) (int, *pkg.ErrorMessage) {
	self.t.Log("Handler-Add", sess, n)
	return n + 1, nil
}

func TestService(t *testing.T) {
	ser := tcp.NewServer("", PORT)
	serv := NewService("test", ser)
	serv.SetEncoding(Encoding)

	handler := NewHandler(t)
	serv.RegistHanlder(handler)

	err := serv.Start()
	assert.Nil(t, err, "servive.Start() error: %v", err)

	cli := tcp.NewClient()
	client := NewServiceClient(cli)
	client.SetEncoding(Encoding)

	err = client.Connect("", PORT)
	assert.Nil(t, err, "servive.Start() error: %v", err)

	//on push
	client.AddListener("test.pushstring", func(line string) {
		t.Log(line)
		assert.Equal(t, "Service: string", line)
	})

	//notify
	err = client.Notify("test.handler.notifystring", "string")
	assert.Nil(t, err, "client.Notify() error: %v", err)

	//request
	err = client.Request("test.handler.requestint", 100, func(result int) {
		assert.Equal(t, 101, result)
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "can't not come to here")
	})
	time.Sleep(1 * time.Second)
	err = serv.Stop()
	assert.Nil(t, err, "serv.Stop() error: %v", err)

	assert.Equal(t, true, handler.Started)
	assert.Equal(t, true, handler.Stopped)
}
