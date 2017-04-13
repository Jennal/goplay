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

import "testing"
import "fmt"
import "reflect"

func TestRouter(t *testing.T) {
	r := NewRouter("gate")
	r.Add(r)
	t.Log(len(r.data))
	for k, v := range r.data {
		t.Logf("%v, %v, %v", k, v.caller, v.method)
	}
}

type test struct {
}

func (tt *test) Add(a int, b float32) float32 {
	fmt.Println("test.Add", a, b)
	return float32(a) + b
}

func TestMethod(t *testing.T) {
	r := NewRouter("test")
	r.Add(&test{})
	m := r.Get("test.test.add")
	result := m.Call(1, float32(2.0))
	t.Log(result...)
	arg := m.NewArg(0)
	t.Log(arg, reflect.TypeOf(arg))
}
