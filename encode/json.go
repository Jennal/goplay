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

//Package encode is Encoder and Decoder for pkg
package encode

import "encoding/json"

//Json encoder/decoder
type Json struct {
}

func (self Json) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (self Json) Unmarshal(data []byte, content interface{}) error {
	return json.Unmarshal(data, content)
}

func NewJson() EncodeDecoder {
	return &Base{
		child: Json{},
	}
}
