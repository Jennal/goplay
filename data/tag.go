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

package data

import (
	"sync"
)

//TagContainer can be set to a object, then check if contains a tag
type TagContainer interface {
	Contains(names ...string) bool
	Add(names ...string)
	Remove(names ...string)
}

type TagContainerImpl struct {
	sync.Mutex
	Tags map[string]bool
}

func NewTagContainer() TagContainer {
	return TagContainerImpl{
		Tags: make(map[string]bool),
	}
}

func (t TagContainerImpl) Contains(names ...string) bool {
	t.Lock()
	defer t.Unlock()

	for _, name := range names {
		_, ok := t.Tags[name]
		if !ok {
			return false
		}
	}

	return true
}

func (t TagContainerImpl) Add(names ...string) {
	t.Lock()
	defer t.Unlock()

	for _, name := range names {
		t.Tags[name] = true
	}
}

func (t TagContainerImpl) Remove(names ...string) {
	t.Lock()
	defer t.Unlock()

	for _, name := range names {
		delete(t.Tags, name)
	}
}
