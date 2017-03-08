package filter

import (
	"github.com/jennal/goplay/pkg"
	"github.com/jennal/goplay/session"
)

type IFilter interface {
	OnNewClient(*session.Session) bool                 /* return false to ignore */
	OnRecv(*session.Session, *pkg.Header, []byte) bool /* return false to ignore */
}
