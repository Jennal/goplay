package service

import (
	"github.com/jennal/goplay/event"
)

type Method struct {
	*event.Method
}

func NewMethod(m interface{}) *Method {
	return &Method{
		Method: event.NewMethod(nil, m),
	}
}
