package resolve

import (
	"fmt"
	"net/http"
	"os"
)

type Engine struct {
	server    *http.Server
	router    *Router
	templates map[string]Template

	middlewares []Handler
	controllers []Controller
}

func Setup() *Engine {
	return &Engine{
		server:      &http.Server{},
		controllers: make([]Controller, 0),
		templates:   make(map[string]Template),
		router:      newRouter(),
	}
}

func (e *Engine) RegisterTemplates(path, extension string, recursive bool) {
	LogMessage("Initialization of html-templates has been started.")

	files, err := os.ReadDir(path)
	if err != nil {
		ErrorMessage(fmt.Sprintf("Failed to load templates from %s", path), -1)
	}

	for _, file := range files {
		if file.IsDir() && recursive {
			e.RegisterTemplates(path+"/"+file.Name(), extension, recursive)
		}

		if file.Name()[len(file.Name())-len(extension):] == extension {
			e.templates[file.Name()] = NewTemplate(file.Name(), path+"/"+file.Name())
		}
	}

	LogMessage("Finished. Templates count: " + fmt.Sprint(len(e.templates)))
}

func (e *Engine) RegisterMiddlewares(middlewares ...Handler) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e *Engine) RegisterControllers(c ...Controller) {
	e.controllers = append(e.controllers, c...)

	for _, c := range c {
		for method, routes := range c.routes {
			for _, route := range routes {
				err := validatePath(route.Path)
				if err != false {
					continue
				}

				route.Path = c.path + route.Path
				route.PathTokens = getPathTokens(route.Path)

				var allHandlers []Handler
				allHandlers = append(allHandlers, e.middlewares...)
				allHandlers = append(allHandlers, route.Handlers...)

				e.router.registerRoute(method, route.Path, allHandlers...)
			}
		}
	}
}

func (e *Engine) Run(port string) {
	go func() {
		e.server.Handler = e
		e.server.Addr = port
		err := e.server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	LogMessage("Server started on port: " + port)
	LogMessage("Press CTRL+C to stop the server...")

	select {}
}

func (e *Engine) Stop() {
	err := e.server.Close()
	if err != nil {
		return
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers, params := e.router.matchRoute(r.Method, r.URL.Path)

	if handlers == nil {
		http.NotFound(w, r)
		return
	}

	ctx := &Ctx{
		app:       e,
		w:         w,
		r:         r,
		handlers:  handlers,
		index:     -1,
		locals:    make(Map),
		binds:     make(Map),
		params:    params,
		templates: e.templates,
	}

	ctx.Next()
}
