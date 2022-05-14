package kid

import (
	"fmt"
	"net/http"
)

// ErrorHandlerFunc defines a function to process return errors or panic errors from handlers.
type ErrorHandlerFunc func(*Ctx, error) error

var errorLogger = NewLogger("HTTP Error")

// DefaultErrorHandler that process return errors or panic errors from handlers.
var DefaultErrorHandler ErrorHandlerFunc = func(c *Ctx, err error) error {
	errorLogger.Error(c, fmt.Sprintf("%s %s", c.Method(), c.Url().RequestURI()), nil, err)
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
