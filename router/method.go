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

package router

import "reflect"

type Method struct {
	caller interface{}
	method reflect.Method
}

func NewMethod(caller interface{}, method reflect.Method) *Method {
	return &Method{
		caller: caller,
		method: method,
	}
}

func (m *Method) Call(args ...interface{}) []interface{} {
	vals := []reflect.Value{
		reflect.ValueOf(m.caller),
	}

	for _, v := range args {
		vals = append(vals, reflect.ValueOf(v))
	}

	result := []interface{}{}
	res := m.method.Func.Call(vals)

	for _, r := range res {
		result = append(result, r.Interface())
	}

	return result
}

func (m *Method) NewArg(i int) interface{} {
	if m.method.Type.NumIn() <= i {
		return nil
	}

	argType := m.method.Type.In(i)
	return reflect.New(argType).Interface()
}
