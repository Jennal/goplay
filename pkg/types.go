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

type Status byte

const (
	STAT_OK                Status = 0x00
	STAT_ERR                      = 0x90
	STAT_ERR_WRONG_PARAMS         = 0x91
	STAT_ERR_DECODE_FAILED        = 0x92
	STAT_ERR_TIMEOUT              = 0x93
	STAT_ERR_EMPTY_RESULT         = 0x94
)

type PackageType byte

const (
	PKG_NOTIFY             PackageType = 0x00
	PKG_REQUEST                        = 0x01
	PKG_RESPONSE                       = 0x02
	PKG_PUSH                           = 0x03
	PKG_HEARTBEAT                      = 0x04
	PKG_HEARTBEAT_RESPONSE             = 0x05

	PKG_RPC_NOTIFY   = 0x06
	PKG_RPC_REQUEST  = 0x07
	PKG_RPC_RESPONSE = 0x08
	PKG_RPC_PUSH     = 0x09
)

type EncodingType byte

const (
	ENCODING_NONE     EncodingType = 0x00
	ENCODING_GOB                   = 0x01
	ENCODING_JSON                  = 0x02
	ENCODING_BSON                  = 0x03
	ENCODING_PROTOBUF              = 0x04
)

type PackageIDType byte
type PackageSizeType uint16
