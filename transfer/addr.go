package transfer

import "net"
import "strings"
import "strconv"

type Addr interface {
	net.Addr
	IP() string
	Port() int
}

type addr struct {
	net.Addr
}

func NewAddr(a net.Addr) Addr {
	return &addr{a}
}

func (a *addr) IP() string {
	if a == nil {
		return ""
	}

	str := a.String()
	idx := strings.LastIndex(str, ":")
	return str[:idx]
}

func (a *addr) Port() int {
	if a == nil {
		return 0
	}

	str := a.String()
	idx := strings.LastIndex(str, ":")
	port, _ := strconv.Atoi(str[idx+1:])
	return port
}
