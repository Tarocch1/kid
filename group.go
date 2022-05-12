package kid

import "net/http"

type group struct {
	prefix string
	kid    *Kid
}

// Group creates a router group.
func (g *group) Group(prefix string) *group {
	group := &group{
		prefix: g.prefix + prefix,
		kid:    g.kid,
	}
	return group
}

func (g *group) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	g.kid.router.addRoute(method, pattern, handler)
}

// Head adds a head router.
func (g *group) Head(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodHead, pattern, handler)
}

// Get adds a get router.
func (g *group) Get(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, pattern, handler)
}

// Delete adds a delete router.
func (g *group) Delete(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodDelete, pattern, handler)
}

// Post adds a post router.
func (g *group) Post(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, pattern, handler)
}

// Put adds a put router.
func (g *group) Put(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPut, pattern, handler)
}

// Patch adds a patch router.
func (g *group) Patch(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPatch, pattern, handler)
}

// Use adds a middleware.
func (g *group) Use(middlewares ...HandlerFunc) {
	g.kid.router.addMiddleware(g.prefix, middlewares...)
}
