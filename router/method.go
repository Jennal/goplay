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

import (
	"fmt"
	"reflect"

	"github.com/jennal/goplay/aop"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
)

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

func (m *Method) Call(sess *session.Session, header *pkg.Header, data []byte) (result []interface{}, err error) {
	result = nil
	err = nil

	if m.NumIn() == 2 {
		aop.Recover(func() {
			result = m.CallArgs(sess)
		}, func(e interface{}) {
			err = fmt.Errorf("router.Method.Call[NumIn=2] err: %#v", e)
			log.RecoverErrorf("%v", err)
		})
	} else if m.NumIn() == 3 {
		val := m.NewArg(2)
		// fmt.Printf("Service.callRouteFunc: %#v => %v\n", val, reflect.TypeOf(val))
		decoder := encode.GetEncodeDecoder(header.Encoding)
		err = decoder.Unmarshal(data, val)
		if err != nil {
			return nil, log.NewErrorf("Service.callRouteFunc decoder.Unmarshal failed: %v", err)
		}
		// fmt.Printf("Service.callRouteFunc: %#v => %v\n", val, reflect.TypeOf(val))

		var arg2 interface{} = val
		if m.method.Type.In(2).Kind() != reflect.Ptr {
			arg2 = helpers.GetValueFromPtr(val)
		}

		aop.Recover(func() {
			result = m.CallArgs(sess, arg2)
		}, func(e interface{}) {
			err = fmt.Errorf("router.Method.Call[NumIn=3] err: %#v", e)
			log.RecoverErrorf("%v", err)
		})
	} else if m.NumIn() == 4 {
		aop.Recover(func() {
			result = m.CallArgs(sess, header, data)
		}, func(e interface{}) {
			err = fmt.Errorf("router.Method.Call[NumIn=4] err: %#v", e)
			log.RecoverErrorf("%v", err)
		})
	}

	if result != nil {
		return
	}

	err = log.NewErrorf("Method.Call can't come to here, must be something wrong, m.NumIn() = %v", m.NumIn())
	return result, err
}

func (m *Method) CallArgs(args ...interface{}) []interface{} {
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
	if argType.Kind() == reflect.Ptr {
		argType = argType.Elem()
	}
	return reflect.New(argType).Interface()
}

func (m *Method) NumIn() int {
	return m.method.Type.NumIn()
}
