package router

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jennal/goplay/handler"
	"github.com/jennal/goplay/session"
)

var (
	TYPE_IHANDLER reflect.Type = reflect.TypeOf((*handler.IHandler)(nil)).Elem()
	TYPE_SESSION  reflect.Type = reflect.TypeOf(session.NewSession(nil))
	TYPE_ERROR                 = reflect.TypeOf((*handler.HandlerError)(nil))
)

type Router struct {
	serverName string
	data       map[string]*Method
}

func NewRouter(serverName string) *Router {
	return &Router{
		serverName: serverName,
		data:       make(map[string]*Method),
	}
}

func (r *Router) Add(obj interface{}) {
	tp := reflect.TypeOf(obj)
	// fmt.Println(tp.NumMethod())
	for i := 0; i < tp.NumMethod(); i++ {
		method := tp.Method(i)
		if !isValidMethod(method) {
			// fmt.Println("1", isValidMethod(method), method)
			continue
		}

		path := r.getPath(tp, method)
		if _, ok := r.data[path]; ok {
			//TODO:Error
			fmt.Println("Router.Add: error: path already exists!", path)
			continue
		}

		r.data[path] = NewMethod(obj, method)
	}
	// fmt.Printf("Router: %#v\n", r.data)
}

func (r *Router) Get(path string) *Method {
	if m, ok := r.data[path]; ok {
		return m
	}

	return nil
}

func (r *Router) getPath(t reflect.Type, m reflect.Method) string {
	return strings.ToLower(fmt.Sprintf("%s.%s.%s",
		r.serverName,
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
	 * func (*handler.IHandler) Method(*session.Session, interface{}) (interface{}, *handler.HandlerError)
	 */

	/* Args: *handler.IHandler, *session.Session, interface{} */
	if m.Type.NumIn() != 3 {
		// fmt.Println("isValidMethod-1")
		return false
	}

	/* Returns: interface{}, *handler.HandlerError */
	if m.Type.NumOut() > 2 {
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

	/* valid return */
	if m.Type.NumOut() == 1 {
		if retType := m.Type.Out(0); retType != TYPE_ERROR {
			// fmt.Println("isValidMethod-5")
			return false
		}

		return true
	}

	if retType := m.Type.Out(1); retType != TYPE_ERROR {
		// fmt.Println("isValidMethod-6")
		return false
	}

	return true
}
