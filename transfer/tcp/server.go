package tcp

import (
	"fmt"
	"net"

	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/transfer"
)

type server struct {
	*event.Event

	host      string
	port      int
	clients   []transfer.IClient
	isStarted bool

	listener net.Listener
}

func NewServer(host string, port int) transfer.IServer {
	return &server{
		Event:     event.NewEvent(),
		host:      host,
		port:      port,
		clients:   []transfer.IClient{},
		isStarted: false,
	}
}

func (serv *server) RegistDelegate(delegate transfer.IServerDelegate) {
	serv.On(transfer.EVENT_SERVER_STARTED, delegate, delegate.OnStarted)
	serv.On(transfer.EVENT_SERVER_STOPPED, delegate, delegate.OnStopped)
	serv.On(transfer.EVENT_SERVER_ERROR, delegate, delegate.OnError)
	serv.On(transfer.EVENT_SERVER_NEW_CLIENT, delegate, delegate.OnNewClient)
}

func (serv *server) UnregistDelegate(delegate transfer.IServerDelegate) {
	serv.Off(transfer.EVENT_SERVER_STARTED, delegate)
	serv.Off(transfer.EVENT_SERVER_STOPPED, delegate)
	serv.Off(transfer.EVENT_SERVER_ERROR, delegate)
	serv.Off(transfer.EVENT_SERVER_NEW_CLIENT, delegate)
}

func (serv *server) GetClients() []transfer.IClient {
	return serv.clients
}

func (serv *server) Start() error {
	if serv.isStarted {
		return nil
	}

	host := fmt.Sprintf("%s:%d", serv.host, serv.port)
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	serv.isStarted = true

	serv.listener = ln
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				serv.Emit(transfer.EVENT_SERVER_ERROR, err)
				serv.Stop()
			}

			client := NewClientWithConnect(conn)
			client.On(transfer.EVENT_CLIENT_DISCONNECTED, serv, func(cli transfer.IClient) {
				//remove from serv.clients
				for i := len(serv.clients) - 1; i >= 0; i-- {
					if serv.clients[i] == cli {
						serv.clients = append(serv.clients[:i], serv.clients[i+1:]...)
						break
					}
				}
			})
			serv.clients = append(serv.clients, client)
			serv.Emit(transfer.EVENT_SERVER_NEW_CLIENT, client)
		}
	}()

	defer serv.Emit(transfer.EVENT_SERVER_STARTED)
	return nil
}

func (serv *server) Stop() error {
	if !serv.isStarted {
		return nil
	}

	defer serv.Emit(transfer.EVENT_SERVER_STOPPED)
	serv.isStarted = false
	return serv.listener.Close()
}
