package kid

import (
	"fmt"
	"net/http"
)

type handler struct {
	kid *Kid
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newCtx(w, req)
	handlerFunc, params, _ := h.kid.router.getRoute(c.Method(), c.Url().Path)
	c.params = params
	middlewares := h.kid.router.getMiddlewares(c.Url().Path)

	handlers := append(middlewares, func(c *Ctx) error {
		if handlerFunc != nil {
			return handlerFunc(c)
		} else {
			return NewError(
				http.StatusNotFound,
				fmt.Sprintf("404 Not Found: %s %s", c.Method(), c.Url().RequestURI()),
				nil,
			)
		}
	})

	c.handlers = handlers
	err := c.Next()
	if err != nil {
		h.kid.config.ErrorHandler(c, err)
	}
}
