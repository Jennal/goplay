package router

import "reflect"
import "fmt"
import "strings"

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
	fmt.Println(tp.NumMethod())
	for i := 0; i < tp.NumMethod(); i++ {
		method := tp.Method(i)
		if !isValidMethod(method) {
			fmt.Println("1", isValidMethod(method), method)
			continue
		}

		path := r.getPath(tp, method)
		if _, ok := r.data[path]; ok {
			//TODO:Error
			continue
		}

		r.data[path] = NewMethod(obj, method)
	}
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
	//TODO:
	return true
}
