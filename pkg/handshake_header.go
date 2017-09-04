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

type HandShakeClientData struct {
	ClientType    string
	ClientVersion string
	DictMd5       string
}

type HandShakeResponse struct {
	ServerVersion string
	Now           string
	HeartBeatRate int
	Routes        map[string]RouteIndex
}

func NewHandShakeHeader(e EncodingType) *Header {
	return &Header{
		Type:         PKG_HAND_SHAKE,
		Encoding:     e,
		ID:           0,
		Status:       STAT_OK,
		ContentSize:  0,
		Route:        "",
		RouteEncoded: ROUTE_INDEX_NONE,
	}
}

func NewHandShakeResponseHeader(h *Header) *Header {
	newHB := &Header{}
	*newHB = *h
	newHB.Type = PKG_HAND_SHAKE_RESPONSE
	return newHB
}
