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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeConvert(t *testing.T) {
	var u32 uint32 = 1024
	buf, err := GetBytes(u32)
	assert.Nil(t, err)
	t.Log(buf)
	u32New, err := ToUInt32(buf)
	assert.Equal(t, u32, u32New)

	var u16 uint16 = 1024
	buf, err = GetBytes(u16)
	assert.Nil(t, err)
	t.Log(buf)
	u16New, err := ToUInt16(buf)
	assert.Nil(t, err)
	assert.Equal(t, u16, u16New)
}
