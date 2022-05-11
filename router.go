package kid

import "strings"

type routerTreeNode struct {
	pattern    string
	part       string
	children   []*routerTreeNode
	isPartWild bool
	isFullWild bool
	handler    HandlerFunc
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

func (t *routerTree) insert(parts []string, handler HandlerFunc) {
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
		}
		cur.children = append(cur.children, next)
		cur = next
	}
	cur.handler = handler
}

func (t *routerTree) search(parts []string) *routerTreeNode {
	cur := t.root
	for _, part := range parts {
		next := cur.matchChild(part)
		if next == nil {
			return nil
		}
		if next.isFullWild {
			return next
		}
		cur = next
	}
	return cur
}

type router struct {
	trees map[string]*routerTree
}

func newRouter() *router {
	return &router{trees: make(map[string]*routerTree)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := toParts(pattern)
	if _, ok := r.trees[method]; !ok {
		r.trees[method] = &routerTree{root: &routerTreeNode{}}
	}
	r.trees[method].insert(parts, handler)
}

func (r *router) getRoute(method string, path string) (HandlerFunc, map[string]string) {
	pathParts := toParts(path)
	tree, ok := r.trees[method]
	if !ok {
		return nil, nil
	}

	n := tree.search(pathParts)
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
		return n.handler, params
	}
	return nil, nil
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
