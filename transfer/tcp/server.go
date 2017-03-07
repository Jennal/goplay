package tcp

import (
	"fmt"
	"net"

	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/transfer"
)

type server struct {
	*event.Event

	host    string
	port    int
	clients []transfer.IClient

	listener net.Listener
}

func NewServer(host string, port int) transfer.IServer {
	return &server{
		Event:   event.NewEvent(),
		host:    host,
		port:    port,
		clients: []transfer.IClient{},
	}
}

func (serv *server) RegistHandler(handler transfer.IServerHandler) {
	serv.On(transfer.EVENT_SERVER_STARTED, handler, handler.OnStarted)
	serv.On(transfer.EVENT_SERVER_STOPPED, handler, handler.OnStopped)
	serv.On(transfer.EVENT_SERVER_ERROR, handler, handler.OnError)
	serv.On(transfer.EVENT_SERVER_NEW_CLIENT, handler, handler.OnNewClient)
}

func (serv *server) UnregistHandler(handler transfer.IServerHandler) {
	serv.Off(transfer.EVENT_SERVER_STARTED, handler)
	serv.Off(transfer.EVENT_SERVER_STOPPED, handler)
	serv.Off(transfer.EVENT_SERVER_ERROR, handler)
	serv.Off(transfer.EVENT_SERVER_NEW_CLIENT, handler)
}

func (serv *server) GetClients() []transfer.IClient {
	return serv.clients
}

func (serv *server) Start() error {
	host := fmt.Sprintf("%s:%d", serv.host, serv.port)
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	serv.listener = ln
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				serv.Emit(transfer.EVENT_SERVER_ERROR, err)
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

	serv.Emit(transfer.EVENT_SERVER_STARTED)
	return nil
}

func (serv *server) Stop() error {
	serv.Emit(transfer.EVENT_SERVER_STOPPED)
	return serv.listener.Close()
}
