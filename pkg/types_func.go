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

//Package pkg defines data structures pass through network
package pkg

func (t PackageType) HasResponse() bool {
	return t == PKG_HAND_SHAKE ||
		t == PKG_HEARTBEAT ||
		t == PKG_REQUEST ||
		t == PKG_RPC_REQUEST // ||
	// t == PKG_NOTIFY ||
	// t == PKG_RPC_NOTIFY
}

func (t PackageType) ToResponse() PackageType {
	switch t {
	case PKG_HAND_SHAKE:
		return PKG_HAND_SHAKE_RESPONSE
	case PKG_HEARTBEAT:
		return PKG_HEARTBEAT_RESPONSE
	case PKG_REQUEST:
		return PKG_RESPONSE
	/* RPC */
	case PKG_RPC_HAND_SHAKE:
		return PKG_RPC_HAND_SHAKE_RESPONSE
	case PKG_RPC_REQUEST:
		return PKG_RPC_RESPONSE
		// case PKG_NOTIFY:
		// 	return PKG_PUSH
		// case PKG_RPC_NOTIFY:
		// 	return PKG_RPC_PUSH
	}

	return t
}

func (t PackageType) IsRPC() bool {
	return t&PKG_RPC == PKG_RPC
}
