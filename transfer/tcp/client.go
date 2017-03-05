package tcp

import (
	"errors"
	"fmt"
	"net"

	"github.com/jennal/goplay/handler/pkg"
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

	fmt.Println("Write:", header, data, buffer)

	// fmt.Println("write-0")
	_, err = client.Write(buffer)
	// fmt.Println("write-1")
	return err
}

func (client *client) Recv(header *pkg.Header, data interface{}) error {
	var buffer = make([]byte, protocol.HEADER_STATIC_SIZE)
	_, err := client.Read(buffer)
	if err != nil {
		return err
	}
	fmt.Println("Header:", err, buffer)

	routeBuf := make([]byte, 1)
	_, err = client.Read(routeBuf)
	if err != nil {
		return err
	}
	buffer = append(buffer, routeBuf...)
	/* heartbeat/heartbeat_response has no route */
	if routeBuf[0] > 0 {
		routeBuf = make([]byte, routeBuf[0])
		_, err = client.Read(routeBuf)
		if err != nil {
			return err
		}

		buffer = append(buffer, routeBuf...)
	}

	h, _, err := protocol.UnMarshalHeader(buffer)
	fmt.Println(h)
	if err != nil {
		return err
	}
	*header = *h

	decoder := protocol.GetEncodeDecoder(header.Encoding)
	buffer = make([]byte, header.ContentSize)
	_, err = client.Read(buffer)
	if err != nil {
		return err
	}

	return decoder.UnmarshalContent(buffer, data)
}
