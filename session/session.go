package session

import (
	"github.com/jennal/goplay/data"
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

func (s *Session) Request(route string, data interface{}, callback interface{}) error {
	// header := s.NewHeader(pkg.PKG_REQUEST, pkg.ENCODING_JSON, route)
	//TODO:

	return nil
}

func (s *Session) Notify(route string, data interface{}) error {
	//TODO:
	return nil
}

func (s *Session) Push(route string, data interface{}) error {
	return s.Notify(route, data)
}

func (s *Session) AddListener(route string, callback interface{}) {
	//TODO:
}
