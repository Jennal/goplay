package handler

import "github.com/jennal/goplay/session"

type IHandler interface {
	OnStarted()
	OnStopped()
	OnNewClient(*session.Session)
}
