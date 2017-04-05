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

//Package event provide a nodejs like event emmiter
package event

import "sync"

type IEvent interface {
	On(string, interface{}, EventFunc)
	Off(string, interface{})
	Once(string, interface{}, EventFunc)
	Emit(string, ...interface{})
}

type EventFunc interface{}

type Event struct {
	sync.Mutex
	cbs map[string][]*Method
}

func NewEvent() *Event {
	return &Event{
		cbs: make(map[string][]*Method),
	}
}

func (self Event) On(name string, ins interface{}, cb EventFunc) {
	self.Lock()
	defer self.Unlock()

	list, ok := self.cbs[name]
	if !ok {
		list = []*Method{
			NewMethod(ins, cb),
		}
	} else {
		list = append(list, NewMethod(ins, cb))
	}
	self.cbs[name] = list
}

func (self Event) Off(name string, ins interface{}) {
	self.Lock()
	defer self.Unlock()

	list, ok := self.cbs[name]
	if !ok {
		return
	}

	for i := len(list) - 1; i >= 0; i-- {
		item := list[i]
		if item.caller != ins {
			continue
		}

		if i+1 > len(list) {
			list = list[:i]
		} else {
			list = append(list[:i], list[i+1:]...)
		}
	}
	self.cbs[name] = list
}

func (self Event) Once(name string, ins interface{}, cb EventFunc) {
	self.On(name, ins, func(args ...interface{}) {
		NewMethod(ins, cb).Call(args...)
		self.Off(name, ins)
	})
}

func (self Event) Emit(name string, args ...interface{}) {
	self.Lock()
	list, ok := self.cbs[name]
	if !ok {
		self.Unlock()
		return
	}

	copyList := make([]*Method, len(list))
	copy(copyList, list)
	self.Unlock()

	for _, item := range copyList {
		// log.Logf("Emit: %v | %#v", name, item)
		item.Call(args...)
	}
}

// func (self Event) EmitParalle(onfinish func(), name string, args ...interface{}) {
// 	list, ok := self.cbs[name]
// 	if !ok {
// 		return
// 	}

// 	funcs := make([]func(complete chan bool), len(list))
// 	for i, item := range list {
// 		funcs[i] = func(complete chan bool) {
// 			item.Call(args...)
// 			complete <- true
// 		}
// 	}

// 	go func() {
// 		aop.Parallel(funcs...)
// 		onfinish()
// 	}()
// }

// func (self Event) EmitSequence(onfinish func(), name string, args ...interface{}) {
// 	list, ok := self.cbs[name]
// 	if !ok {
// 		return
// 	}

// 	funcs := make([]func(next chan bool, exit chan bool), len(list))
// 	for i, item := range list {
// 		funcs[i] = func(next chan bool, exit chan bool) {
// 			item.Call(args...)
// 			next <- true
// 		}
// 	}

// 	go func() {
// 		aop.Sequence(funcs...)
// 		onfinish()
// 	}()
// }
