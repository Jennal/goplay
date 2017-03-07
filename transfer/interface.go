package transfer

import (
	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/pkg"
)

type IClient interface {
	event.IEvent

	IsConnected() bool
	Connect(host string, port int) error
	Disconnect() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)

	Send(*pkg.Header, []byte) error
	Recv() (*pkg.Header, []byte, error)
}

type IServer interface {
	event.IEvent

	RegistHandler(handler IServerHandler)
	UnregistHandler(handler IServerHandler)

	GetClients() []IClient
	Start() error
	Stop() error
}

type IServerHandler interface {
	OnStarted()
	OnError(err error)
	OnStopped()

	OnNewClient(client IClient)
}
