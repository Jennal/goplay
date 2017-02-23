package transfer

import (
	"github.com/jennal/goplay/handler/pkg"
)

type Client interface {
	IsConnected() bool
	Connect(host string, port int) error
	Disconnect() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)

	Send(*pkg.Header, interface{}) error
	Recv(*pkg.Header, interface{}) error
}

type Server interface {
	GetClients() []Client
	Start() error
	Stop() error
}

type ServerHandler interface {
	OnStarted()
	OnError(err error)
	OnStopped()

	OnNewClient(client Client)
}
