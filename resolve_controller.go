package resolve

type Controller struct {
	path   string
	routes map[string][]Route
}

func NewController(path string) Controller {
	return Controller{
		path:   path,
		routes: make(map[string][]Route),
	}
}

func (c *Controller) registerRoute(method, path string, handlers ...Handler) {
	err := validatePath(path)
	if err != false {
		return
	}

	c.routes[method] = append(c.routes[method], Route{
		Path:     path,
		Handlers: handlers,
	})
}

func (c *Controller) Get(path string, handlers ...Handler) {
	c.registerRoute("GET", path, handlers...)
}

func (c *Controller) Post(path string, handlers ...Handler) {
	c.registerRoute("POST", path, handlers...)
}

func (c *Controller) Put(path string, handlers ...Handler) {
	c.registerRoute("PUT", path, handlers...)
}

func (c *Controller) Delete(path string, handlers ...Handler) {
	c.registerRoute("DELETE", path, handlers...)
}

func (c *Controller) Patch(path string, handlers ...Handler) {
	c.registerRoute("PATCH", path, handlers...)
}

func (c *Controller) Options(path string, handlers ...Handler) {
	c.registerRoute("OPTIONS", path, handlers...)
}
