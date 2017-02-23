package service

import (
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/transfer"
)

type Service interface {
}

type service struct {
	server transfer.Server
	router router.Router
}
