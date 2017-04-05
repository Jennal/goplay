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

package channel

import (
	"sync"
)

type ChannelManager struct {
	sync.Mutex
	channels map[string]*Channel
}

var (
	staticChannelManager *ChannelManager
)

func GetChannelManager() *ChannelManager {
	if staticChannelManager == nil {
		staticChannelManager = &ChannelManager{
			channels: make(map[string]*Channel),
		}
	}

	return staticChannelManager
}

func (cm *ChannelManager) ChannelNames() []string {
	cm.Lock()
	defer cm.Unlock()

	result := make([]string, len(cm.channels))
	for k := range cm.channels {
		result = append(result, k)
	}

	return result
}

func (cm *ChannelManager) Channels() map[string]*Channel {
	cm.Lock()
	defer cm.Unlock()

	result := make(map[string]*Channel)
	for k, v := range cm.channels {
		result[k] = v
	}

	return result
}

func (cm *ChannelManager) Create(name string) *Channel {
	ch := NewChannel(name)

	cm.Lock()
	defer cm.Unlock()
	cm.channels[name] = ch

	return ch
}

func (cm *ChannelManager) Get(name string) *Channel {
	cm.Lock()
	defer cm.Unlock()

	if ch, ok := cm.channels[name]; ok {
		return ch
	}

	return nil
}
