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

//Package session stores client data in memory
package session

import (
	"math"

	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
)

var IDGen = helpers.NewIDGen(math.MaxUint32)

type Session struct {
	transfer.IClient

	ID       uint32
	ClientID uint32
	Encoding pkg.EncodingType
	Encoder  encode.EncodeDecoder
}

func NewSession(cli transfer.IClient) *Session {
	return &Session{
		IClient:  cli,
		ID:       0,
		ClientID: 0,
		Encoding: defaults.Encoding,
		Encoder:  encode.GetEncodeDecoder(defaults.Encoding),
	}
}

func (s *Session) Bind(id uint32) {
	s.ID = id
}

func (s *Session) BindClientID(id uint32) {
	s.ClientID = id
}

func (s *Session) SetEncoding(e pkg.EncodingType) error {
	if encoder := encode.GetEncodeDecoder(e); encoder != nil {
		s.Encoding = e
		s.Encoder = encoder
		return nil
	}

	return log.NewErrorf("can't find encoder with: %v", e)
}

func (s *Session) Push(route string, data interface{}) error {
	buf, err := s.Encoder.Marshal(data)
	if err != nil {
		return err
	}
	return s.PushRaw(route, buf)
}

func (s *Session) PushRaw(route string, data []byte) error {
	header := s.NewHeader(pkg.PKG_PUSH, s.Encoding, route)
	if s.ClientID != 0 {
		// log.Log("Push: clientId = ", s.ClientID)
		header = pkg.NewRpcHeader(header, s.ClientID)
	}
	return s.Send(header, data)
}

func (s *Session) Broadcast(route string, data []byte) error {
	header := s.NewHeader(pkg.PKG_RPC_BROADCAST, s.Encoding, route)
	if s.ClientID != 0 {
		// log.Log("Push: clientId = ", s.ClientID)
		header = pkg.NewRpcHeader(header, s.ClientID)
	}
	return s.Send(header, data)
}
