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
	"testing"

	"github.com/jennal/goplay/log"
)

func TestRecover(t *testing.T) {
	Func1()
}

func Func1() {
	Func2()
}

func Func2() {
	Recover(func() {
		panic("Hello Error")
	}, func(err interface{}) {
		log.RecoverErrorf("error: %#v", err)
	})
}
