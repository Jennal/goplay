package handler

import (
	"fmt"

	"github.com/jennal/goplay/pkg"
)

type HandlerError struct {
	error
	Code    pkg.Status
	Message string
}

func (self HandlerError) Error() string {
	return fmt.Sprintf("[%d]: %s", self.Code, self.Message)
}
