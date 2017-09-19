package common

import "net"

type ErrorTimeout struct {
	msg string
}

func NewErrorTimeout(msg string) *net.OpError {
	return &net.OpError{
		Err: &ErrorTimeout{msg: msg},
	}
}

func (et *ErrorTimeout) Timeout() bool {
	return true
}

func (et *ErrorTimeout) Error() string {
	return et.msg
}
