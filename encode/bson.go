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

package encode

import "gopkg.in/mgo.v2/bson"

type Bson struct {
}

func (self Bson) Marshal(obj interface{}) ([]byte, error) {
	return bson.Marshal(obj)
}

func (self Bson) Unmarshal(data []byte, content interface{}) error {
	return bson.Unmarshal(data, content)
}

func NewBson() EncodeDecoder {
	return &Base{
		child: Bson{},
	}
}
