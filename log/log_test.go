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

package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	Log(1, 2, 3)
	Logf("%v | %v | %v", 1, 2, 3)
	Trace(1, 2, 3)
	Tracef("%v | %v | %v", 1, 2, 3)

	Error(errors.New("err-1"))
	Errorf("%v-%v", "err", 2)
	err := NewError("err", 3)
	assert.Error(t, err)
	err = NewErrorf("%v-%v", "err", 3)
	assert.Error(t, err)
	RecoverErrorf("%v-%v", "err", 4)
}
