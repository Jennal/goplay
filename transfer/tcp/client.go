package tcp

import (
	"errors"
	"fmt"
	"net"

	"github.com/jennal/goplay/event"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
)

var (
	ERR_ALREADY_CONNECTED = errors.New("already connected")
)

type client struct {
	*event.Event
	conn        net.Conn
	isConnected bool
}

func NewClientWithConnect(conn net.Conn) transfer.IClient {
	return &client{
		Event:       event.NewEvent(),
		conn:        conn,
		isConnected: true,
	}
}

func NewClient() transfer.IClient {
	return &client{
		Event:       event.NewEvent(),
		isConnected: false,
	}
}

func (client *client) IsConnected() bool {
	return client.isConnected
}

func (client *client) Connect(host string, port int) error {
	if client.isConnected {
		return ERR_ALREADY_CONNECTED
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
	defer client.Emit(transfer.EVENT_CLIENT_DISCONNECTED, client)
	return client.conn.Close()
}

func (client *client) Read(buf []byte) (int, error) {
	return client.conn.Read(buf)
}

func (client *client) Write(buf []byte) (int, error) {
	return client.conn.Write(buf)
}

func (client *client) Send(header *pkg.Header, data []byte) error {
	header.ContentSize = pkg.PackageSizeType(len(data))
	headerBuffer, err := header.Marshal()
	if err != nil {
		return err
	}
	buffer := append(headerBuffer, data...)
	// fmt.Println("Write:", header, data, buffer)

	_, err = client.Write(buffer)

	return err
}

func (client *client) Recv() (*pkg.Header, []byte, error) {
	header := &pkg.Header{}
	_, err := pkg.ReadHeader(client, header)
	if err != nil {
		return nil, nil, err
	}

	if header.ContentSize > 0 {
		buffer := make([]byte, header.ContentSize)
		_, err := client.Read(buffer)
		if err != nil {
			return nil, nil, err
		}

		// fmt.Println("Recv body:", buffer)
		return header, buffer, err
	}

	return header, nil, nil
}
