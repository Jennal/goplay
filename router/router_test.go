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

package router

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
)

type test struct {
	handler.IHandler
}

func (tt *test) OnStarted()                   {}
func (tt *test) OnStopped()                   {}
func (tt *test) OnNewClient(*session.Session) {}

func (tt *test) Add(sess *session.Session, a int) (int, *pkg.ErrorMessage) {
	fmt.Println("test.Add", a)
	return a + 1, nil
}

func (tt *test) Get(sess *session.Session) (int, *pkg.ErrorMessage) {
	fmt.Println("test.Add")
	return 100, nil
}

func TestMethodNumIn(t *testing.T) {
	caller := &test{}
	r := NewRouter()
	r.Add("gate", caller)

	for k, v := range r.data {
		t.Logf("%v, %v, %v, %v", k, v.caller, v.method, v.NumIn())
	}
}

func TestRouter(t *testing.T) {
	caller := &test{}
	r := NewRouter()
	r.Add("gate", caller)
	t.Log(len(r.data))
	assert.Equal(t, 2, len(r.data))
	for k, v := range r.data {
		t.Logf("%v, %v, %v", k, v.caller, v.method)
	}

	v, ok := r.data["gate.test.add"]
	assert.True(t, ok)
	assert.Equal(t, caller, v.caller)
}

func TestMethod(t *testing.T) {
	r := NewRouter()
	r.Add("test", &test{})
	m := r.Get("test.test.add")
	assert.NotNil(t, m)
	result := m.Call(session.NewSession(nil), 1)
	t.Log(result...)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 2, result[0])
	assert.Nil(t, result[1])
	arg := m.NewArg(0)
	t.Log(arg, reflect.TypeOf(arg))

	m = r.Get("test.test.get")
	assert.NotNil(t, m)
	result = m.Call(session.NewSession(nil))
	t.Log(result...)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 100, result[0])
	assert.Nil(t, result[1])
	arg = m.NewArg(0)
	t.Log(arg, reflect.TypeOf(arg))
}
