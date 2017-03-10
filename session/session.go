package session

import (
	"github.com/jennal/goplay/data"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
)

type Session struct {
	transfer.IClient
	*data.Map

	ID int
}

func NewSession(cli transfer.IClient) *Session {
	return &Session{
		IClient: cli,
		ID:      0,
		Map:     data.NewMap(),
	}
}

func (s *Session) Bind(id int) {
	s.ID = id
}

func (s *Session) Push(encoding pkg.EncodingType, route string, data interface{}) error {
	header := s.NewHeader(pkg.PKG_NOTIFY, encoding, route)
	encoder := encode.GetEncodeDecoder(encoding)
	buf, err := encoder.Marshal(data)
	if err != nil {
		return err
	}
	return s.Send(header, buf)
}
