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
