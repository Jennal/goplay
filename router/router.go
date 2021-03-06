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

//Package router resolve network package route and find the right method to call
package router

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
)

var (
	TYPE_IHANDLER = reflect.TypeOf((*handler.IHandler)(nil)).Elem()
	TYPE_SESSION  = reflect.TypeOf((*session.Session)(nil))
	TYPE_ERROR    = reflect.TypeOf((*pkg.ErrorMessage)(nil))
	TYPE_HEADER   = reflect.TypeOf((*pkg.Header)(nil))
	TYPE_BYTES    = reflect.TypeOf(([]byte)(nil))
)

type Router struct {
	data map[string]*Method
}

func NewRouter() *Router {
	return &Router{
		data: make(map[string]*Method),
	}
}

// func (r *Router) calcRoutesMap() map[string]pkg.RouteIndex {
// 	arr := r.GetKeys()
// 	m := make(map[string]pkg.RouteIndex)

// 	for i, str := range arr {
// 		m[str] = pkg.RouteIndex(i + 1)
// 	}

// 	return m
// }

func (r *Router) GetKeys() []string {
	arr := []string{}

	for k := range r.data {
		arr = append(arr, k)
	}

	return arr
}

func (r *Router) Add(prefix string, obj interface{}) {
	tp := reflect.TypeOf(obj)
	// fmt.Println(tp.NumMethod())
	for i := 0; i < tp.NumMethod(); i++ {
		method := tp.Method(i)
		if !isValidMethod(method) {
			// fmt.Println("1", isValidMethod(method), method)
			continue
		}

		path := r.getPath(prefix, tp, method)
		if _, ok := r.data[path]; ok {
			log.Errorf("Router.Add: error: path(%v) already exists!", path)
			continue
		}

		r.data[path] = NewMethod(obj, method)
		pkg.DefaultHandShake().AddRoutesMap(path)
	}
	// fmt.Printf("Router: %#v\n", r.data)
}

func (r *Router) Get(path string) *Method {
	if m, ok := r.data[path]; ok {
		return m
	}

	return nil
}

func (r *Router) getPath(prefix string, t reflect.Type, m reflect.Method) string {
	return strings.ToLower(fmt.Sprintf("%s.%s.%s",
		prefix,
		getStructName(t.String()),
		m.Name,
	))
}

func getStructName(name string) string {
	arr := strings.Split(name, ".")
	arr = arr[len(arr)-1:]
	return arr[0]
}

func isValidMethod(m reflect.Method) bool {
	/*
	 * valid method:
	 * == Request
	 * func (*handler.IHandler) Method(*session.Session, *pkg.Header, []byte) (interface{}, *pkg.ErrorMessage)
	 * func (*handler.IHandler) Method(*session.Session, interface{}) (interface{}, *pkg.ErrorMessage)
	 * func (*handler.IHandler) Method(*session.Session) (interface{}, *pkg.ErrorMessage)
	 * == Notify
	 * func (*handler.IHandler) Method(*session.Session, *pkg.Header, []byte) *pkg.ErrorMessage
	 * func (*handler.IHandler) Method(*session.Session, interface{}) *pkg.ErrorMessage
	 * func (*handler.IHandler) Method(*session.Session) *pkg.ErrorMessage
	 * interface of In(2) / In(1) and Out(0) should not be reflect.Ptr
	 */

	/* Args:
	 * NumIn == 2 : *handler.IHandler, *session.Session
	 * NumIn == 3 : *handler.IHandler, *session.Session, interface{}
	 * NumIn == 4 : *handler.IHandler, *session.Session, *pkg.Header, []byte
	 */
	if m.Type.NumIn() != 2 && m.Type.NumIn() != 3 && m.Type.NumIn() != 4 {
		// fmt.Println("isValidMethod-1")
		return false
	}

	/* Returns: *pkg.ErrorMessage | (interface{}, *pkg.ErrorMessage) */
	if m.Type.NumOut() != 1 && m.Type.NumOut() != 2 {
		// fmt.Println("isValidMethod-2")
		return false
	}

	/* valid args */
	if selfType := m.Type.In(0); !selfType.Implements(TYPE_IHANDLER) {
		// fmt.Println("isValidMethod-3")
		return false
	}

	if sessType := m.Type.In(1); sessType.Kind() != reflect.Ptr || sessType != TYPE_SESSION {
		// fmt.Println("isValidMethod-4")
		return false
	}

	// allow Ptr for In(2)
	// if m.Type.NumIn() == 3 && m.Type.In(2).Kind() == reflect.Ptr {
	// 	// fmt.Println("isValidMethod-5")
	// 	return false
	// }

	if m.Type.NumIn() == 4 &&
		(m.Type.In(2) != TYPE_HEADER ||
			m.Type.In(3) != TYPE_BYTES) {
		return false
	}

	/* valid return */
	if m.Type.NumOut() == 1 {
		if retType := m.Type.Out(0); retType != TYPE_ERROR {
			// fmt.Println("isValidMethod-6")
			return false
		}

		return true
	}

	// allow Ptr for In(2)
	// if m.Type.Out(0).Kind() == reflect.Ptr {
	// 	// fmt.Println("isValidMethod-7")
	// 	return false
	// }

	if retType := m.Type.Out(1); retType != TYPE_ERROR {
		// fmt.Println("isValidMethod-8")
		return false
	}

	return true
}
