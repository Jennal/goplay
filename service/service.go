package service

import (
	"github.com/jennal/goplay/router"
	"github.com/jennal/goplay/transfer"
)

type Service struct {
	server transfer.Server
	router router.Router
}
