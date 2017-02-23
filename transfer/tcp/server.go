package tcp

import (
	"fmt"
	"net"

	"github.com/jennal/goplay/transfer"
)

type server struct {
	host    string
	port    int
	clients []transfer.Client

	delegate transfer.ServerHandler
	listener net.Listener
}

func NewServer(host string, port int, delegate transfer.ServerHandler) transfer.Server {
	return &server{
		host:     host,
		port:     port,
		clients:  []transfer.Client{},
		delegate: delegate,
	}
}

func (serv *server) GetClients() []transfer.Client {
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
				serv.delegate.OnError(err)
			}

			client := NewClientWithConnect(conn)
			serv.clients = append(serv.clients, client)
			serv.delegate.OnNewClient(client)
		}
	}()

	defer serv.delegate.OnStarted()
	return nil
}

func (serv *server) Stop() error {
	defer serv.delegate.OnStopped()
	return serv.listener.Close()
}
