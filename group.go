package kid

import "net/http"

type group struct {
	prefix      string
	middlewares []HandlerFunc
	kid         *Kid
}

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
	g.kid.router.AddRoute(method, pattern, handler)
}

func (g *group) HEAD(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodHead, pattern, handler)
}

func (g *group) GET(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, pattern, handler)
}

func (g *group) DELETE(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodDelete, pattern, handler)
}

func (g *group) POST(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, pattern, handler)
}

func (g *group) PUT(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPut, pattern, handler)
}

func (g *group) PATCH(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPatch, pattern, handler)
}

func (g *group) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
