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

package helpers

type IDGen struct {
	maxID  uint32
	nextID uint32
}

func NewIDGen(max uint32) *IDGen {
	return &IDGen{
		maxID:  max,
		nextID: 0,
	}
}

func (self *IDGen) NextID() uint32 {
	if self.nextID == self.maxID {
		defer func() { self.nextID = 0 }()
	} else {
		defer func() { self.nextID++ }()
	}

	return self.nextID
}
