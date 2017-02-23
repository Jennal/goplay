package tcp

import (
	"errors"
	"fmt"
	"net"

	"github.com/jennal/goplay/handler/pkg"
	"github.com/jennal/goplay/helpers"
	"github.com/jennal/goplay/protocol"
	"github.com/jennal/goplay/transfer"
)

type client struct {
	conn        net.Conn
	isConnected bool
}

func NewClientWithConnect(conn net.Conn) transfer.Client {
	return &client{
		conn:        conn,
		isConnected: true,
	}
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

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	client.conn = conn
	client.isConnected = true

	return nil
}

func (client *client) Disconnect() error {
	return client.conn.Close()
}

func (client *client) Read(buf []byte) (int, error) {
	return client.conn.Read(buf)
}

func (client *client) Write(buf []byte) (int, error) {
	return client.conn.Write(buf)
}

func (client *client) Send(header *pkg.Header, data interface{}) error {
	encoder := protocol.GetEncodeDecoder(header.Encoding)
	buffer, err := encoder.Marshal(header, data)
	if err != nil {
		return err
	}

	size := len(buffer)
	fmt.Println("Send size:", size)
	sizeBuf, err := helpers.UInt32(size).GetBytes()
	if err != nil {
		return err
	}
	_, err = client.Write(sizeBuf)
	if err != nil {
		return err
	}

	_, err = client.Write(buffer)
	return err
}

func (client *client) Recv(header *pkg.Header, data interface{}) error {
	var buffer = make([]byte, 4)
	_, err := client.Read(buffer)
	if err != nil {
		return err
	}

	size, err := helpers.Bytes(buffer).ToInt()
	fmt.Println("Recv size:", size)
	if err != nil {
		return err
	}

	buffer = make([]byte, size)
	_, err = client.Read(buffer)
	if err != nil {
		return err
	}

	h, n, err := protocol.UnMarshalHeader(buffer)
	if err != nil {
		return err
	}

	*header = *h
	decoder := protocol.GetEncodeDecoder(header.Encoding)
	return decoder.UnmarshalContent(buffer[n:], data)
}
