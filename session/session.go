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
