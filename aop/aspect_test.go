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
	"fmt"
	"testing"
	"time"
)

func TestAop(t *testing.T) {
	NewAspect().
		Retry(3).
		Delay(10 * time.Second).
		Repeat(5).
		Do(func() {
			fmt.Println("Test")
		})
}

func TestSequence(t *testing.T) {
	Sequence(func(next chan bool, exit chan bool) int {
		t.Log(1, "=>", 1)
		next <- true
		return 1
	}, func(next chan bool, exit chan bool, a int) int {
		t.Log(2, "=>", a)
		next <- true
		return a + 1
	}, func(next chan bool, exit chan bool, a int) int {
		t.Log(3, "=>", a)
		next <- true
		return a + 1
	}, func(next chan bool, exit chan bool, a int) int {
		t.Log(4, "=>", a)
		exit <- true
		return a + 1
	}, func(next chan bool, exit chan bool, a int) {
		t.Log(5, "=>", a)
		next <- true
	})
}

func TestParallel(t *testing.T) {
	Parallel(func(complete chan bool) {
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
}
