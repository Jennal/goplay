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

package event

import "testing"
import "fmt"

type TestIns struct {
	Name string
}

func (self TestIns) Callback(str string) {
	fmt.Println(self.Name, str)
}

type TestInsMulti struct {
	Name string
}

func (self TestInsMulti) Callback(t *testing.T, a, b int) {
	fmt.Println(self.Name, a, b)
}

func TestEvent(t *testing.T) {
	evt := NewEvent()
	ins := &TestIns{"Name-1"}
	insOnce := &TestIns{"Name-2"}

	m := NewMethod(ins, ins.Callback)
	m.Call("1")

	evt.On("test", ins, ins.Callback)
	evt.Once("test", insOnce, insOnce.Callback)
	evt.On("test", ins, ins.Callback)
	fmt.Println("========================")
	evt.Emit("test", "1")
	fmt.Println("========================")
	evt.Emit("test", "2")
	fmt.Println("========================")

	evt.Off("test", ins)
	evt.Emit("test", "3")
	fmt.Println("========================")
}

func TestMethod(t *testing.T) {
	ins := &TestInsMulti{"Name-1"}
	m := NewMethod(ins, ins.Callback)
	m.Call(t, 1+1, 2)

	func(args ...interface{}) {
		m.Call(args...)
	}(t, 1, 1)
}

func TestPPPMethod(t *testing.T) {
	f := func(args ...interface{}) {
		t.Log(args...)
	}

	m := NewMethod(nil, f)
	t.Logf("%v, %v", m.method.Type().In(0), m.method.Type().NumIn())
}
