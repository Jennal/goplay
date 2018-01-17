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

package session

import (
	"testing"

	"math"

	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/transfer/tcp"
	"github.com/stretchr/testify/assert"
)

var idgen = helpers.NewIDGen(math.MaxInt32)

func TestIdGen(t *testing.T) {
	assert.Equal(t, uint32(0), idgen.NextID())
	assert.Equal(t, uint32(1), idgen.NextID())
}

func TestEvent(t *testing.T) {
	sess := NewSession(tcp.NewClient())
	count := 0
	sess.On("1", sess, func() {
		count++
	})
	sess.Emit("1")
	assert.Equal(t, 1, count)

	sess.Off("1", sess)

	sess.Once("1", sess, func() {
		count++
	})
	sess.Emit("1")
	assert.Equal(t, 2, count)
	sess.Emit("1")
	assert.Equal(t, 2, count)
}

func TestPushCache(t *testing.T) {
	sess := NewSession(tcp.NewClient())
	sess.PushCache("1", 1)
	// sess.PushCache("2", 2)

	t.Log(sess.PopAllCaches())
	t.Log(sess.PopAllCaches())
}
