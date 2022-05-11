package kid

import (
	"net/http"
)

type ErrorHandlerFunc func(*Ctx, error) error

// DefaultErrorHandler that process return errors or panic errors from handlers.
var DefaultErrorHandler ErrorHandlerFunc = func(c *Ctx, err error) error {
	if e, ok := err.(*Error); ok {
		return c.Status(e.Status).String(e.Message)
	} else {
		return c.Status(http.StatusInternalServerError).String(err.Error())
	}
}

type Error struct {
	// Http status
	Status int

	// Message to show
	Message string
}

// NewError creates *kid.Error.
func NewError(status int, message string) *Error {
	return &Error{
		Status:  status,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
