package kid

import (
	"net/http"
)

type ErrorHandlerFunc func(*Ctx, error) error

type Error struct {
	Status  int
	Message string
}

func NewError(status int, message string) *Error {
	return &Error{
		Status:  status,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func defaultErrorHandler(c *Ctx, err error) {
	if e, ok := err.(*Error); ok {
		c.Status(e.Status).String(e.Message)
	} else {
		c.Status(http.StatusInternalServerError).String(err.Error())
	}
}
