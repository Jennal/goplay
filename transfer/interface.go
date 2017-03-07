package transfer

import (
	"github.com/jennal/goplay/pkg"
)

type Client interface {
	IsConnected() bool
	Connect(host string, port int) error
	Disconnect() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)

	Send(*pkg.Header, []byte) error
	Recv() (*pkg.Header, []byte, error)
}

type Server interface {
	SetupHandler(handler ServerHandler)
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
