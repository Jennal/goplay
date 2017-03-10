package transfer

import (
	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer/base"
)

type IClient interface {
	event.IEvent
	base.IHeaderCreator

	RegistDelegate(delegate IClientDelegate)
	UnregistDelegate(delegate IClientDelegate)

	IsConnected() bool
	Connect(host string, port int) error
	Disconnect() error
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)

	Send(*pkg.Header, []byte) error
	Recv() (*pkg.Header, []byte, error)
}

type IClientDelegate interface {
	OnConnected(IClient)
	OnDisconnected(IClient)
}

type IServer interface {
	event.IEvent

	RegistDelegate(delegate IServerDelegate)
	UnregistDelegate(delegate IServerDelegate)

	GetClients() []IClient
	Start() error
	Stop() error
}

type IServerDelegate interface {
	OnStarted()
	OnError(err error)
	OnStopped()

	OnNewClient(client IClient)
}
