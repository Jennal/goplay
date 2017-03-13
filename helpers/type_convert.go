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
	"bytes"
	"encoding/binary"

	"github.com/jennal/goplay/log"
)

func GetBytes(i interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, i)

	return buffer.Bytes(), err
}

func ToUInt32(buffer []byte) (uint32, error) {
	if len(buffer) < 4 {
		return 0, log.NewError("length of buffer < 4")
	}

	var i uint32 = 0
	r := bytes.NewReader(buffer[:4])
	err := binary.Read(r, binary.BigEndian, &i)

	return i, err
}

func ToUInt16(buffer []byte) (uint16, error) {
	if len(buffer) < 2 {
		return 0, log.NewError("length of buffer < 2")
	}

	var i uint16 = 0
	r := bytes.NewReader(buffer[:2])
	err := binary.Read(r, binary.BigEndian, &i)

	return i, err
}
