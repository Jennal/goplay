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

//Package aop is Aspect Oriented Programming.
package aop

import (
	"time"
)

/*
Aspect can do calls like chain

	NewAspect().
	    Retry(3).
	    Delay(10 * time.Second).
	    Repeat(5).
	    Do(func() {
	        fmt.Println("Test")
	})

Means print Test line 5 times, after 10 seconds, if panic happen retry 3 times.
*/
type Aspect interface {
	Recover(err func(interface{})) Aspect
	Retry(times int) Aspect
	Repeat(times int) Aspect
	Delay(d time.Duration) Aspect
	Do(work func())
}

//NewAspect is constructor of Aspect
func NewAspect() Aspect {
	return &aspect{func(do func()) {
		do()
	}}
}

type aspect struct {
	chain func(func())
}

func (a *aspect) Join(work func(func())) Aspect {
	// fmt.Println("Joining", work)
	oldChain := a.chain
	a.chain = func(do func()) {
		// fmt.Println("chaining", work)
		oldChain(func() {
			work(do)
		})
	}

	return a
}

func (a *aspect) Recover(onErr func(interface{})) Aspect {
	a.Join(func(do func()) {
		Recover(do, onErr)
	})

	return a
}

func (a *aspect) Retry(times int) Aspect {
	// fmt.Println("Joining Retry", times)
	a.Join(func(do func()) {
		Retry(do, times)
	})

	return a
}

func (a *aspect) Repeat(times int) Aspect {
	// fmt.Println("Joining Repeat", times)
	a.Join(func(do func()) {
		Repeat(do, times)
	})

	return a
}

func (a *aspect) Delay(d time.Duration) Aspect {
	a.Join(func(do func()) {
		Delay(do, d)
	})

	return a
}

func (a *aspect) Do(work func()) {
	// fmt.Println("Doing", work)
	if a.chain != nil {
		// fmt.Printf("chain %#v\n", a.chain)
		a.chain(work)
	} else {
		// fmt.Println("work", work)
		work()
	}
}
