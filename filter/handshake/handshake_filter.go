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
	"github.com/jennal/goplay/transfer"

	"github.com/jennal/goplay/encode"

	"github.com/jennal/goplay/consts"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/session"
)

type HandShakeFilter struct {
	router *router.Router
}

func NewHandShakeFilter(r *router.Router) *HandShakeFilter {
	return &HandShakeFilter{
		router: r,
	}
}

func SendHandShake(client transfer.IClient, e pkg.EncodingType, encoder encode.EncodeDecoder, md5 string) error {
	buffer, err := encoder.Marshal(pkg.HandShakeClientData{
		ClientType:    consts.ClientType,
		ClientVersion: consts.Version,
		DictMd5:       md5,
	})
	if err != nil {
		return err
	}
	client.Send(pkg.NewHandShakeHeader(e), buffer)
	return nil
}

func (self *HandShakeFilter) OnNewClient(sess *session.Session) bool { /* return false to ignore */
	return true
}

func (self *HandShakeFilter) OnRecv(sess *session.Session, header *pkg.Header, data []byte) bool { /* return false to ignore */
	if header.Type != pkg.PKG_HAND_SHAKE {
		return true
	}

	encoder := encode.GetEncodeDecoder(header.Encoding)

	decodeData := &pkg.HandShakeClientData{}
	err := encoder.Unmarshal(data, decodeData)
	if err != nil {
		log.Error(err)
		return true
	}

	routeMap := pkg.HandShakeInstance.RoutesMap()
	respData := pkg.HandShakeResponse{
		ServerVersion: consts.Version,
		Now:           time.Now(),
		HeartBeatRate: consts.HeartBeatRate,
	}

	//check md5
	encodeRouteMap, err := encoder.Marshal(routeMap)
	if err != nil {
		log.Error(err)
		return true
	}

	if decodeData.DictMd5 != helpers.Md5(encodeRouteMap) {
		respData.Routes = routeMap
	}

	encodeRespData, err := encoder.Marshal(respData)
	if err != nil {
		log.Error(err)
		return true
	}

	header.Type = pkg.PKG_HAND_SHAKE_RESPONSE
	sess.Send(header, encodeRespData)

	return false
}
