package kid

import (
	"net/http"
)

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(*Ctx) error

// Kid application.
type Kid struct {
	*group
	router *router
	config Config
}

// New creates a kid app.
func New(config ...Config) *Kid {
	kid := &Kid{
		router: newRouter(),
		config: Config{},
	}
	kid.group = &group{kid: kid}
	if len(config) > 0 {
		kid.config = config[0]
	}
	setDefaultConfig(kid)
	return kid
}

// Listen starts server at addr.
func (k *Kid) Listen(addr string) (err error) {
	handler := &handler{
		kid: k,
	}
	return http.ListenAndServe(addr, handler)
}

// ListenTLS starts server at addr with https.
func (k *Kid) ListenTLS(addr string, certFile string, keyFile string) (err error) {
	handler := &handler{
		kid: k,
	}
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}
