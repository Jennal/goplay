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

package pkg

func NewHeartBeatHeader(idGen *IDGen) *Header {
	return NewHeader(PKG_HEARTBEAT, ENCODING_NONE, idGen, "")
}

func NewHeartBeatResponseHeader(h *Header) *Header {
	newHB := &Header{}
	*newHB = *h
	newHB.Type = PKG_HEARTBEAT_RESPONSE
	return newHB
}