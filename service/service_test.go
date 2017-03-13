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

type CustomMessage struct {
	Name string
}

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
	self.t.Log("Handler-NotifyString", sess, line)
	assert.Equal(self.t, "string", line)
	sess.Push("test.pushstring", "Service: "+line)
	return nil
}

func (self *Handler) RequestInt(sess *session.Session, n int) (int, *pkg.ErrorMessage) {
	self.t.Log("Handler-RequestInt", sess, n)
	return n + 1, nil
}

func (self *Handler) RequestString(sess *session.Session, n string) (string, *pkg.ErrorMessage) {
	self.t.Log("Handler-RequestString", sess, n)
	return n + "-1", nil
}

func (self *Handler) RequestArray(sess *session.Session, n []string) ([]string, *pkg.ErrorMessage) {
	self.t.Log("Handler-RequestArray", sess, n)
	return append(n, "Hello to Client"), nil
}

func (self *Handler) RequestMap(sess *session.Session, n map[string]int) (map[string]int, *pkg.ErrorMessage) {
	self.t.Log("Handler-RequestMap", sess, n)
	n["Hi Client"] = 100
	return n, nil
}

func (self *Handler) RequestObj(sess *session.Session, n CustomMessage) (CustomMessage, *pkg.ErrorMessage) {
	self.t.Log("Handler-RequestObj", sess, n)
	n.Name = "From Service"
	return n, nil
}

func (self *Handler) RequestFail(sess *session.Session, n int) (int, *pkg.ErrorMessage) {
	self.t.Log("Handler-RequestFail", sess, n)
	return 0, pkg.NewErrorMessage(pkg.STAT_ERR_WRONG_PARAMS, "Test Error")
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
		t.Log("[test.pushstring] Recv => ", line)
		assert.Equal(t, "Service: string", line)
	})

	//notify
	err = client.Notify("test.handler.notifystring", "string")
	assert.Nil(t, err, "client.Notify() error: %v", err)

	//request
	//int
	err = client.Request("test.handler.requestint", 100, func(result int) {
		t.Log("[test.handler.requestint] Recv => ", result)
		assert.Equal(t, 101, result)
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "can't not come to here")
	})
	assert.Nil(t, err, "client.Request() error: %v", err)

	//string
	err = client.Request("test.handler.requeststring", "100", func(result string) {
		t.Log("[test.handler.requeststring] Recv => ", result)
		assert.Equal(t, "100-1", result)
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "can't not come to here")
	})
	assert.Nil(t, err, "client.Request() error: %v", err)

	//array
	err = client.Request("test.handler.requestarray", []string{"Hello to Service"}, func(result []string) {
		t.Log("[test.handler.requestarray] Recv => ", result)
		assert.Equal(t, []string{"Hello to Service", "Hello to Client"}, result)
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "can't not come to here")
	})
	assert.Nil(t, err, "client.Request() error: %v", err)

	//map
	err = client.Request("test.handler.requestmap", map[string]int{
		"Hello to Service": 10,
	}, func(result map[string]int) {
		t.Log("[test.handler.requestmap] Recv => ", result)
		assert.Equal(t, map[string]int{
			"Hello to Service": 10,
			"Hi Client":        100,
		}, result)
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "can't not come to here")
	})
	assert.Nil(t, err, "client.Request() error: %v", err)

	//object
	err = client.Request("test.handler.requestobj", CustomMessage{
		Name: "Hello to Service",
	}, func(result CustomMessage) {
		t.Log("[test.handler.requestobj] Recv => ", result)
		assert.Equal(t, CustomMessage{
			Name: "From Service",
		}, result)
	}, func(err *pkg.ErrorMessage) {
		assert.True(t, false, "can't not come to here")
	})
	assert.Nil(t, err, "client.Request() error: %v", err)

	//error
	err = client.Request("test.handler.requestfail", 0, func(result int) {
		assert.True(t, false, "can't not come to here")
	}, func(err *pkg.ErrorMessage) {
		t.Log("Recv Error:", err)
		assert.Equal(t, pkg.NewErrorMessage(pkg.STAT_ERR_WRONG_PARAMS, "Test Error"), err)
	})
	assert.Nil(t, err, "client.Request() error: %v", err)

	time.Sleep(1 * time.Second)
	err = serv.Stop()
	assert.Nil(t, err, "serv.Stop() error: %v", err)

	assert.Equal(t, true, handler.Started)
	assert.Equal(t, true, handler.Stopped)
}
