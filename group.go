package kid

import "net/http"

type group struct {
	prefix      string
	middlewares []HandlerFunc
	kid         *Kid
}

// Group creates a router group.
func (g *group) Group(prefix string) *group {
	group := &group{
		prefix: g.prefix + prefix,
		kid:    g.kid,
	}
	g.kid.groups = append(g.kid.groups, group)
	return group
}

func (g *group) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	g.kid.router.addRoute(method, pattern, handler)
}

// HEAD adds a head router.
func (g *group) HEAD(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodHead, pattern, handler)
}

// GET add a get router.
func (g *group) GET(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, pattern, handler)
}

// DELETE add a delete router.
func (g *group) DELETE(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodDelete, pattern, handler)
}

// POST add a post router.
func (g *group) POST(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, pattern, handler)
}

// PUT add a put router.
func (g *group) PUT(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPut, pattern, handler)
}

// PATCH add a patch router.
func (g *group) PATCH(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPatch, pattern, handler)
}

// Use add a middleware.
func (g *group) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
