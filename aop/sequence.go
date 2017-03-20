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

/*Sequence can make parallel work in sequence, and 2nd func can have input from 1st func's return value, and so on.

Example:
	aop.Sequence(func(next chan bool, exit chan bool) int {
		fmt.Println(1, "=>", 1)
		next <- true
		return 1
	}, func(next chan bool, exit chan bool, a int) int {
		fmt.Println(2, "=>", a)
		next <- true
		return a + 1
	}, func(next chan bool, exit chan bool, a int) int {
		fmt.Println(3, "=>", a)
		next <- true
		return a + 1
	}, func(next chan bool, exit chan bool, a int) int {
		fmt.Println(4, "=>", a)
		exit <- true
		return a + 1
	}, func(next chan bool, exit chan bool, a int) {
		fmt.Println(5, "=>", a)
		next <- true
	})
*/
func Sequence(funcs ...interface{}) {
	var in []reflect.Value
	var nextChan = make(chan bool, 1)
	var exitChan = make(chan bool, 1)
	chanIn := []reflect.Value{
		reflect.ValueOf(nextChan),
		reflect.ValueOf(exitChan),
	}

Loop:
	for _, v := range funcs {
		go func() {
			in = append(chanIn, in...)
			in = reflect.ValueOf(v).Call(in)
		}()

		select {
		case <-nextChan:
			continue
		case <-exitChan:
			break Loop
		}
	}
}
