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
	message := fmt.Sprintf("%s %s", c.Method(), c.Url().RequestURI())
	if e, ok := err.(*Error); ok {
		errorLogger.Error(c, message, map[string]interface{}{
			"data": err.(*Error).Data,
		}, err)
		return c.Status(e.Status).String(e.Message)
	} else {
		errorLogger.Error(c, message, nil, err)
		return c.Status(http.StatusInternalServerError).String(err.Error())
	}
}

type Error struct {
	// Http status
	Status int

	// Message to show
	Message string

	// Extra data
	Data interface{}
}

// NewError creates *kid.Error.
func NewError(status int, message string, data interface{}) *Error {
	return &Error{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func (e *Error) Error() string {
	return e.Message
}
