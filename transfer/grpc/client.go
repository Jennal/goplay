package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/transfer"
)

type client struct {
	conn         *grpc.ClientConn
	streamClient Grpc_StreamClient
	grpcClient   GrpcClient
	isConnected  bool
}

func NewClient() transfer.Client {
	return &client{
		isConnected: false,
	}
}

func (client *client) IsConnected() bool {
	return client.isConnected
}

func (client *client) Connect(host string, port int) error {
	if client.isConnected {
		return errors.New("already connected")
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return err
	}

	client.conn = conn
	client.isConnected = true

	client.grpcClient = NewGrpcClient(client.conn)
	client.streamClient, err = client.grpcClient.Stream(context.Background())
	if err != nil {
		return err
	}

	// go func() {
	// 	client.streamClient.Send(&Package{
	// 		Type:     1,
	// 		Encoding: 2,
	// 		ID:       3,
	// 		Content:  []byte("Hello From Client"),
	// 	})

	// 	for {
	// 		pkg, err := client.streamClient.Recv()
	// 		fmt.Println("Recv:", pkg, err)
	// 	}
	// }()

	return nil
}

func (client *client) Disconnect() error {
	return client.conn.Close()
}

func (client *client) Read(buf []byte) (int, error) {
	return 0, nil
}

func (client *client) Write(buf []byte) (int, error) {
	client.streamClient.Send(&Package{
		Type:     1,
		Encoding: 2,
		ID:       3,
		Content:  buf,
	})
	return 0, nil
}

func (client *client) Send(header *pkg.Header, data interface{}) error {
	// encoder := protocol.GetEncodeDecoder(header.Encoding)
	// buffer, err := encoder.Marshal(header, data)
	// if err != nil {
	// 	return err
	// }

	// var size = len(buffer)
	// fmt.Println("Send size:", size)
	// sizeBuf, err := helpers.UInt32(size).GetBytes()
	// if err != nil {
	// 	return err
	// }
	// _, err = client.Write(sizeBuf)
	// if err != nil {
	// 	return err
	// }

	// _, err = client.Write(buffer)
	// return err
	return nil
}

func (client *client) Recv(header *pkg.Header, data interface{}) error {
	// var buffer = make([]byte, 4)
	// _, err := client.Read(buffer)
	// if err != nil {
	// 	return err
	// }

	// size, err := helpers.Bytes(buffer).ToInt()
	// fmt.Println("Recv size:", size)
	// if err != nil {
	// 	return err
	// }

	// buffer = make([]byte, size)
	// _, err = client.Read(buffer)
	// if err != nil {
	// 	return err
	// }

	// h, n, err := protocol.UnMarshalHeader(buffer)
	// if err != nil {
	// 	return err
	// }

	// *header = *h
	// decoder := protocol.GetEncodeDecoder(header.Encoding)
	// return decoder.UnmarshalContent(buffer[n:], data)
	return nil
}
