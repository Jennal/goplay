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

type Status int8

const (
	STAT_OK Status = iota
	STAT_ERR
	STAT_ERR_WRONG_PARAMS
	STAT_ERR_DECODE_FAILED
	STAT_ERR_TIMEOUT
)

type PackageType byte

const (
	PKG_NOTIFY PackageType = iota
	PKG_REQUEST
	PKG_RESPONSE
	PKG_HEARTBEAT
	PKG_HEARTBEAT_RESPONSE

	PKG_RPC_NOTIFY
	PKG_RPC_REQUEST
	PKG_RPC_RESPONSE
)

type EncodingType byte

const (
	ENCODING_NONE EncodingType = iota
	ENCODING_GOB
	ENCODING_JSON
	ENCODING_BSON
	ENCODING_PROTOBUF
)

type PackageIDType byte
type PackageSizeType uint16
