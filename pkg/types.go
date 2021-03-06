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

//go:generate stringer -type=PackageType
type PackageType byte

const (
	PKG_NOTIFY              PackageType = 0x00
	PKG_REQUEST             PackageType = 0x01
	PKG_RESPONSE            PackageType = 0x02
	PKG_PUSH                PackageType = 0x03
	PKG_HEARTBEAT           PackageType = 0x04
	PKG_HEARTBEAT_RESPONSE  PackageType = 0x05
	PKG_HAND_SHAKE          PackageType = 0x06
	PKG_HAND_SHAKE_RESPONSE PackageType = 0x07

	PKG_RPC = 0x10

	PKG_RPC_NOTIFY              PackageType = PKG_RPC | PKG_NOTIFY
	PKG_RPC_REQUEST             PackageType = PKG_RPC | PKG_REQUEST
	PKG_RPC_RESPONSE            PackageType = PKG_RPC | PKG_RESPONSE
	PKG_RPC_PUSH                PackageType = PKG_RPC | PKG_PUSH
	PKG_RPC_HAND_SHAKE          PackageType = PKG_RPC | PKG_HAND_SHAKE
	PKG_RPC_HAND_SHAKE_RESPONSE PackageType = PKG_RPC | PKG_HAND_SHAKE_RESPONSE
	PKG_RPC_BROADCAST           PackageType = PKG_RPC | 0x08
)

//go:generate stringer -type=EncodingType
type EncodingType byte

const (
	ENCODING_NONE     EncodingType = 0x00
	ENCODING_GOB      EncodingType = 0x01
	ENCODING_JSON     EncodingType = 0x02
	ENCODING_BSON     EncodingType = 0x03
	ENCODING_PROTOBUF EncodingType = 0x04
)

type PackageIDType byte
type PackageSizeType uint16

//go:generate stringer -type=RouteIndex
type RouteIndex uint16

const (
	ROUTE_INDEX_NONE RouteIndex = 0
)

type RouteMap map[string]RouteIndex
