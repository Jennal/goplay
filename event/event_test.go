package event

import "testing"
import "fmt"

type TestIns struct {
	Name string
}

func (self TestIns) Callback(str string) {
	fmt.Println(self.Name, str)
}

func TestEvent(t *testing.T) {
	evt := NewEvent()
	ins := &TestIns{"Name-1"}
	insOnce := &TestIns{"Name-2"}

	m := NewMethod(ins, ins.Callback)
	m.Call("1")

	evt.On("test", ins, ins.Callback)
	evt.On("test", ins, ins.Callback)
	evt.Once("test", insOnce, insOnce.Callback)
	fmt.Println("========================")
	evt.Emit("test", "1")
	fmt.Println("========================")
	evt.Emit("test", "2")
	fmt.Println("========================")

	evt.Off("test", ins)
	evt.Emit("test", "3")
	fmt.Println("========================")
}
