package handler

type IHandler interface {
	OnStarted()
	OnStopped()
	OnNewClient()
}
