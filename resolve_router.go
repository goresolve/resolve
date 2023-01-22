package resolve

import (
	"sync"
)

type Router struct {
	routes map[string][]Route

	mu sync.RWMutex
}

func newRouter() *Router {
	router := &Router{
		routes: make(map[string][]Route),
		mu:     sync.RWMutex{},
	}

	return router
}

type pathToken struct {
	IsParametric    bool
	Value           string
	ParametricValue string
}

type Param struct {
	Value string
	Name  string
}

type Route struct {
	Path       string
	PathTokens []pathToken
	Handlers   []Handler `json:"-"`

	Parametric bool
}

func (r *Router) matchRoute(method string, path string) ([]Handler, []Param) {
	pathTokens := getTokens(path)

	for _, route := range r.routes[method] {
		if len(pathTokens) != len(route.PathTokens) {
			continue
		}

		var params []Param
		match := true
		for i, tokens := range pathTokens {
			if route.PathTokens[i].IsParametric {
				params = append(params, Param{
					Name:  route.PathTokens[i].ParametricValue,
					Value: tokens,
				})
				continue
			}
			if tokens != route.PathTokens[i].Value {
				match = false
				break
			}
		}

		if match {
			return route.Handlers, params
		}
	}

	return nil, nil
}

func (r *Router) registerRoute(method string, path string, handler ...Handler) {
	validatePath(path)

	r.routes[method] = append(r.routes[method], Route{
		Path:       path,
		PathTokens: getPathTokens(path),
		Handlers:   handler,
		Parametric: isParametric(path),
	})
}
