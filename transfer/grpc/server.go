package grpc

import (
	"fmt"
	"net"

	"github.com/jennal/goplay/transfer"
	grpc1 "google.golang.org/grpc"
)

type server struct {
	host    string
	port    int
	clients []transfer.Client

	delegate   transfer.ServerHandler
	listener   net.Listener
	grpcServer *grpc1.Server
}

func NewServer(host string, port int, delegate transfer.ServerHandler) transfer.Server {
	return &server{
		host:     host,
		port:     port,
		clients:  []transfer.Client{},
		delegate: delegate,
	}
}

func (serv *server) Stream(gss Grpc_StreamServer) error {
	for {
		pkg, err := gss.Recv()
		if err != nil {
			fmt.Println(err)
			fmt.Println(pkg)
			return nil
		}

		fmt.Println(pkg)
		pkg.Content = []byte("Hello from Server")
		gss.Send(pkg)
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
	serv.grpcServer = grpc1.NewServer()
	RegisterGrpcServer(serv.grpcServer, serv)
	go serv.grpcServer.Serve(serv.listener)

	defer serv.delegate.OnStarted()
	return nil
}

func (serv *server) Stop() error {
	defer serv.delegate.OnStopped()
	return serv.listener.Close()
}
