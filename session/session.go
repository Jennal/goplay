package session

import (
	"github.com/jennal/goplay/data"
	"github.com/jennal/goplay/defaults"
	"github.com/jennal/goplay/encode"
	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/transfer"
)

type Session struct {
	transfer.IClient
	*data.Map

	ID       int
	Encoding pkg.EncodingType
	encoder  encode.EncodeDecoder
}

func NewSession(cli transfer.IClient) *Session {
	return &Session{
		IClient:  cli,
		Map:      data.NewMap(),
		ID:       0,
		Encoding: defaults.Encoding,
	}
}

func (s *Session) Bind(id int) {
	s.ID = id
}

func (s *Session) SetEncoding(e pkg.EncodingType) error {
	if encoder := encode.GetEncodeDecoder(e); encoder != nil {
		s.Encoding = e
		s.encoder = encoder
		return nil
	}

	return log.NewErrorf("can't find encoder with: %v", e)
}

func (s *Session) Push(route string, data interface{}) error {
	header := s.NewHeader(pkg.PKG_NOTIFY, s.Encoding, route)
	encoder := encode.GetEncodeDecoder(s.Encoding)
	buf, err := encoder.Marshal(data)
	if err != nil {
		return err
	}
	return s.Send(header, buf)
}
