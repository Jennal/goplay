package router

import "testing"
import "fmt"
import "reflect"

func TestRouter(t *testing.T) {
	r := NewRouter("gate")
	r.Add(r)
	t.Log(len(r.data))
	for k, v := range r.data {
		t.Logf("%v, %v, %v", k, v.caller, v.method)
	}
}

type test struct {
}

func (tt *test) Add(a int, b float32) float32 {
	fmt.Println("test.Add", a, b)
	return float32(a) + b
}

func TestMethod(t *testing.T) {
	r := NewRouter("test")
	r.Add(&test{})
	m := r.Get("test.test.add")
	result := m.Call(1, float32(2.0))
	t.Log(result...)
	arg := m.NewArg(0)
	t.Log(arg, reflect.TypeOf(arg))
}
