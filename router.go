package kid

import (
	"net/http"
	"strings"
)

const middlewaresMethod = "middlewares"

type routerTreeNode struct {
	pattern     string
	part        string
	children    []*routerTreeNode
	isPartWild  bool
	isFullWild  bool
	middlewares []HandlerFunc
	handler     HandlerFunc
}

func (n *routerTreeNode) matchChild(part string) *routerTreeNode {
	for _, child := range n.children {
		if child.part == part || child.isPartWild || child.isFullWild {
			return child
		}
	}
	return nil
}

type routerTree struct {
	root *routerTreeNode
}

func (t *routerTree) insert(parts []string, middleware bool, handler ...HandlerFunc) {
	cur := t.root
	pattern := ""
	for _, part := range parts {
		next := cur.matchChild(part)
		pattern = pattern + "/" + part
		if next == nil {
			next = &routerTreeNode{
				pattern:    pattern,
				part:       part,
				isPartWild: part[0] == ':',
				isFullWild: part[0] == '*',
			}
			cur.children = append(cur.children, next)
		}
		cur = next
	}
	if middleware {
		cur.middlewares = append(cur.middlewares, handler...)
	} else {
		cur.handler = handler[0]
	}
}

func (t *routerTree) search(parts []string) (*routerTreeNode, []*routerTreeNode) {
	nodes := []*routerTreeNode{t.root}
	cur := t.root
	for _, part := range parts {
		next := cur.matchChild(part)
		if next == nil {
			return nil, nodes
		}
		nodes = append(nodes, next)
		if next.isFullWild {
			return next, nodes
		}
		cur = next
	}
	return cur, nodes
}

type router struct {
	trees map[string]*routerTree
}

func newRouter() *router {
	return &router{trees: map[string]*routerTree{
		http.MethodHead:   {root: &routerTreeNode{}},
		http.MethodGet:    {root: &routerTreeNode{}},
		http.MethodDelete: {root: &routerTreeNode{}},
		http.MethodPost:   {root: &routerTreeNode{}},
		http.MethodPut:    {root: &routerTreeNode{}},
		http.MethodPatch:  {root: &routerTreeNode{}},
		middlewaresMethod: {root: &routerTreeNode{}},
	}}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := toParts(pattern)
	r.trees[method].insert(parts, false, handler)
}

func (r *router) getRoute(method string, path string) (HandlerFunc, map[string]string, []*routerTreeNode) {
	pathParts := toParts(path)
	tree, ok := r.trees[method]
	if !ok {
		return nil, nil, nil
	}

	n, ns := tree.search(pathParts)
	if n != nil {
		parts := toParts(n.pattern)
		params := make(map[string]string)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = pathParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(pathParts[index:], "/")
				break
			}
		}
		return n.handler, params, ns
	}
	return nil, nil, ns
}

func (r *router) addMiddleware(pattern string, middlewares ...HandlerFunc) {
	parts := toParts(pattern)
	r.trees[middlewaresMethod].insert(parts, true, middlewares...)
}

func (r *router) getMiddlewares(path string) []HandlerFunc {
	middlewares := make([]HandlerFunc, 0)
	_, _, ns := r.getRoute(middlewaresMethod, path)
	for _, n := range ns {
		middlewares = append(middlewares, n.middlewares...)
	}
	return middlewares
}

func toParts(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
