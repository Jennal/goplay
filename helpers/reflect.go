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

//Package helpers provide helper functions
package helpers

import "reflect"

var (
	TYPE_BYTES = reflect.TypeOf(([]byte)(nil))
)

func GetValueFromPtr(ptr interface{}) interface{} {
	return reflect.ValueOf(ptr).Elem().Interface()
}

func IsBytesType(obj interface{}) bool {
	return reflect.TypeOf(obj) == TYPE_BYTES
}
