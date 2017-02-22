package transfer

type Client interface {
	IsConnected() bool
	Connect(host string, port int) error
	Recv() ([]byte, error)
	Send([]byte) error
}
