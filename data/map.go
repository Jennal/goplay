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

//Package data is a thread safe memory storage of data
package data

import (
	"sync"
)

//Map is a key value storage, key should be string, value can be any type
type Map struct {
	mutex sync.Mutex
	data  map[string]interface{}
}

//NewMap creates new Map struct
func NewMap() *Map {
	return &Map{
		data: make(map[string]interface{}),
	}
}

func (s *Map) Remove(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
}

func (s *Map) Set(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
}

func (s *Map) HasKey(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, has := s.data[key]
	return has
}

func (s *Map) Int(key string) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Int8(key string) int8 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int8)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Int16(key string) int16 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int16)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Int32(key string) int32 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int32)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Int64(key string) int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int64)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Uint(key string) uint {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Uint8(key string) uint8 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint8)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Uint16(key string) uint16 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint16)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Uint32(key string) uint32 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint32)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Uint64(key string) uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint64)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Float32(key string) float32 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float32)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) Float64(key string) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float64)
	if !ok {
		return 0
	}
	return value
}

func (s *Map) String(key string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.data[key]
	if !ok {
		return ""
	}

	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

func (s *Map) Value(key string) interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.data[key]
}

// Retrieve all Map state
func (s *Map) State() map[string]interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.data
}

// Restore Map state after reconnect
func (s *Map) Restore(data map[string]interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = data
}
