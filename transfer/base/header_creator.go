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

//Package base provide common logic to transfer
package base

import "github.com/jennal/goplay/pkg"

type HeaderCreator struct {
	idGen *pkg.IDGen
}

func NewHeaderCreator() *HeaderCreator {
	return &HeaderCreator{
		idGen: pkg.NewIDGen(),
	}
}

func (self *HeaderCreator) NewHeader(t pkg.PackageType, e pkg.EncodingType, r string) *pkg.Header {
	return pkg.NewHeader(t, e, self.idGen, r)
}

func (self *HeaderCreator) NewHeartBeatHeader() *pkg.Header {
	return pkg.NewHeartBeatHeader(self.idGen)
}

func (self *HeaderCreator) NewHeartBeatResponseHeader(h *pkg.Header) *pkg.Header {
	return pkg.NewHeartBeatResponseHeader(h)
}
