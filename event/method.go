package event

import "reflect"

type Method struct {
	caller interface{}
	method reflect.Value
}

func NewMethod(caller interface{}, method interface{}) *Method {
	return &Method{
		caller: caller,
		method: reflect.ValueOf(method),
	}
}

func (m *Method) Call(args ...interface{}) []interface{} {
	vals := []reflect.Value{}

	for _, v := range args {
		vals = append(vals, reflect.ValueOf(v))
	}

	result := []interface{}{}
	res := m.method.Call(vals)

	for _, r := range res {
		result = append(result, r.Interface())
	}

	return result
}
