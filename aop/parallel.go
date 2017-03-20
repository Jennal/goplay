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

package aop

import (
	"reflect"
)

/*Parallel can make works in parallel, and wait for all complete.

Example:
	aop.Parallel(func(complete chan bool) {
		t.Log(1)
		complete <- true
	}, func(complete chan bool) {
		t.Log(2)
		complete <- true
	}, func(complete chan bool) {
		t.Log(3)
		complete <- true
	}, func(complete chan bool) {
		t.Log(4)
		complete <- true
	}, func(complete chan bool) {
		t.Log(5)
		complete <- true
	})
*/
func Parallel(funcs ...interface{}) {
	completeChan := make(chan bool, len(funcs))
	chanIn := []reflect.Value{reflect.ValueOf(completeChan)}

	for _, v := range funcs {
		f := v
		go func() {
			reflect.ValueOf(f).Call(chanIn)
		}()
	}

	for i := 0; i < len(funcs); i++ {
		<-completeChan
	}
}
