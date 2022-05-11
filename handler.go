package kid

import (
	"fmt"
	"net/http"
	"strings"
)

type handler struct {
	kid *Kid
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newCtx(w, req)
	handlerFunc, params := h.kid.router.GetRoute(c.Method(), c.Url().Path)
	c.params = params

	middlewares := make([]HandlerFunc, 0)
	for _, group := range h.kid.groups {
		if strings.HasPrefix(c.Url().Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	handlers := append(middlewares, func(c *Ctx) error {
		if handlerFunc != nil {
			return handlerFunc(c)
		} else {
			return NewError(
				http.StatusNotFound,
				fmt.Sprintf("404 Not Found: %s %s", c.Method(), c.Url().RequestURI()),
			)
		}
	})

	c.handlers = handlers
	err := c.Next()
	if err != nil {
		if h.kid.config.ErrorHandler != nil {
			h.kid.config.ErrorHandler(c, err)
		} else {
			defaultErrorHandler(c, err)
		}
	}
}
