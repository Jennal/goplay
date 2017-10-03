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

//Package handshake is the first service processor that handles route encode
package handshake

import (
	"time"

	"github.com/jennal/goplay/helpers"

	"github.com/jennal/goplay/encode"

	"github.com/jennal/goplay/consts"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
)

type HandShakeFilter struct {
	reconnectTo *pkg.HostPort
}

func NewHandShakeFilter() *HandShakeFilter {
	return &HandShakeFilter{}
}

func GetHandShakeRequest(e pkg.EncodingType, encoder encode.EncodeDecoder, md5 string) (*pkg.Header, []byte, error) {
	buffer, err := encoder.Marshal(&pkg.HandShakeClientData{
		ClientType:    consts.ClientType,
		ClientVersion: consts.Version,
		DictMd5:       md5,
	})
	if err != nil {
		return nil, nil, err
	}

	return pkg.NewHandShakeHeader(e), buffer, nil
}

func (self *HandShakeFilter) SetReconnectTo(data *pkg.HostPort) {
	self.reconnectTo = data
}

func (self *HandShakeFilter) OnNewClient(sess *session.Session) bool { /* return false to ignore */
	return true
}

func (self *HandShakeFilter) OnRecv(sess *session.Session, header *pkg.Header, data []byte) bool { /* return false to ignore */
	if header.Type&^pkg.PKG_RPC != pkg.PKG_HAND_SHAKE {
		return true
	}

	encoder := encode.GetEncodeDecoder(header.Encoding)

	decodeData := &pkg.HandShakeClientData{}
	err := encoder.Unmarshal(data, decodeData)
	if err != nil {
		log.Error(err)
		return true
	}

	var routeMap pkg.RouteMap
	routeMap = pkg.DefaultHandShake().RoutesMap()

	respData := &pkg.HandShakeResponse{
		ServerVersion: consts.Version,
		Now:           time.Now().Format("2006-01-02 15:04:05"),
		HeartBeatRate: consts.HeartBeatRate,
		Routes:        nil,

		IsReconnect: false,
		ReconnectTo: &pkg.HostPort{
			Host: "",
			Port: 0,
		},
	}

	if self.reconnectTo != nil {
		respData.IsReconnect = true
		respData.ReconnectTo = self.reconnectTo
	}

	//check md5
	md5Encoder := encode.GetMd5EncodeDecoder()
	encodeRouteMap, err := md5Encoder.Marshal(routeMap)
	if err != nil {
		log.Error(err)
		return true
	}

	if decodeData.DictMd5 != helpers.Md5(encodeRouteMap) {
		respData.Routes = make(map[string]uint32)
		for k, v := range routeMap {
			respData.Routes[k] = uint32(v)
		}
	}

	encodeRespData, err := encoder.Marshal(respData)
	if err != nil {
		log.Error(err)
		return true
	}

	header.Type = header.Type.ToResponse()
	// log.Logf("HandShake Response:\n\tins(%p) => %#v\n\theader => %#v\n\t%#v | %v", pkg.DefaultHandShake(), pkg.DefaultHandShake(), header, encodeRespData, string(encodeRespData))
	err = sess.Send(header, encodeRespData)
	if err != nil {
		log.Error(err)
		return true
	}

	return false
}
