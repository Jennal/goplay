package pkg

import (
	"fmt"
)

type ErrorMessage struct {
	Code    Status
	Message string
}

func NewErrorMessage(code Status, msg string) *ErrorMessage {
	return &ErrorMessage{
		Code:    code,
		Message: msg,
	}
}

func (self ErrorMessage) Error() string {
	return fmt.Sprintf("[%d]: %s", self.Code, self.Message)
}
