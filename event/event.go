package event

import "sync"

type IEvent interface {
	On(string, interface{}, EventFunc)
	Off(string, interface{})
	Once(string, interface{}, EventFunc)
	Emit(string, ...interface{})
}

type EventFunc interface{}

type Event struct {
	sync.Mutex
	cbs map[string][]*Method
}

func NewEvent() *Event {
	return &Event{
		cbs: make(map[string][]*Method),
	}
}

func (self Event) On(name string, ins interface{}, cb EventFunc) {
	self.Lock()
	defer self.Unlock()

	list, ok := self.cbs[name]
	if !ok {
		list = []*Method{
			NewMethod(ins, cb),
		}
	} else {
		list = append(list, NewMethod(ins, cb))
	}
	self.cbs[name] = list
}

func (self Event) Off(name string, ins interface{}) {
	self.Lock()
	defer self.Unlock()

	list, ok := self.cbs[name]
	if !ok {
		return
	}

	for i := len(list) - 1; i >= 0; i-- {
		item := list[i]
		if item.caller != ins {
			continue
		}

		if i+1 > len(list) {
			list = list[:i]
		} else {
			list = append(list[:i], list[i+1:]...)
		}
	}
	self.cbs[name] = list
}

func (self Event) Once(name string, ins interface{}, cb EventFunc) {
	self.On(name, ins, func(args ...interface{}) {
		NewMethod(ins, cb).Call(args...)
		self.Off(name, ins)
	})
}

func (self Event) Emit(name string, args ...interface{}) {
	list, ok := self.cbs[name]
	if !ok {
		return
	}

	for _, item := range list {
		item.Call(args...)
	}
}
