package kid

import (
	"net/http"
)

type HandlerFunc func(*Ctx) error

type Kid struct {
	*group
	router *router
	groups []*group
	config *Config
}

func New(config *Config) *Kid {
	kid := &Kid{router: newRouter()}
	kid.group = &group{kid: kid}
	kid.groups = []*group{kid.group}
	kid.config = config
	return kid
}

func (k *Kid) Listen(addr string) (err error) {
	handler := &handler{
		kid: k,
	}
	return http.ListenAndServe(addr, handler)
}

func (k *Kid) ListenTLS(addr string, certFile string, keyFile string) (err error) {
	handler := &handler{
		kid: k,
	}
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}
