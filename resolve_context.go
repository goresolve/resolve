package resolve

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Ctx struct {
	app *Engine

	w http.ResponseWriter
	r *http.Request

	templates map[string]Template
	handlers  []Handler
	index     int

	locals Map
	binds  Map

	params []Param
}

func (c *Ctx) App() *Engine {
	return c.app
}

func (c *Ctx) Request() *http.Request {
	return c.r
}

func (c *Ctx) Response() http.ResponseWriter {
	return c.w
}

func (c *Ctx) Next() {
	c.index++

	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Ctx) Bind(binds Map) {
	for key, value := range binds {
		c.binds[key] = value
	}
}

func (c *Ctx) SetLocal(key string, value interface{}) {
	c.locals[key] = value
}

func (c *Ctx) GetLocal(key string) interface{} {
	return c.locals[key]
}

func (c *Ctx) Locals() Map {
	return c.locals
}

func (c *Ctx) Binds() Map {
	return c.binds
}

func (c *Ctx) Param(key string) (string, error) {
	for _, param := range c.params {
		if param.Name == key {
			return param.Value, nil
		}
	}

	ErrorMessage("param not found: "+key, 3)

	return "", errors.New("param not found")
}

func (c *Ctx) HTML(file string, data Map) {
	tmpl, ok := c.templates[file]
	if !ok {
		ErrorMessage("template not found: "+file, 3)
		return
	}

	for k, v := range c.binds {
		data[k] = v
	}

	tmpl.template.Execute(c.Response(), data)
}

func (c *Ctx) Status(status Status) *Ctx {
	c.Response().WriteHeader(int(status))
	return c
}

func (c *Ctx) Redirect(url string) error {
	c.Response().Header().Set("Location", url)
	c.Response().WriteHeader(int(StatusSeeOther))

	return nil
}

func (c *Ctx) RedirectStatus(url string, status Status) error {
	c.Response().Header().Set("Location", url)
	c.Response().WriteHeader(int(status))

	return nil
}

func (c *Ctx) SendString(msg string) error {
	c.Response().Write([]byte(msg))

	return nil
}

func (c *Ctx) Send(msg []byte) error {
	c.Response().Write(msg)

	return nil
}

func (c *Ctx) Form(key, def string) string {
	if val := c.Request().FormValue(key); val != "" {
		return val
	}

	return def
}

func (c *Ctx) FormFile(key string) (multipart.File, error) {
	file, _, err := c.Request().FormFile(key)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (c *Ctx) SaveFile(file multipart.File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}

func (c *Ctx) Query(key, def string) string {
	if val := c.Request().URL.Query().Get(key); val != "" {
		return val
	}

	return def
}

func (c *Ctx) Body() ([]byte, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Ctx) IsGet() bool {
	if c.Request().Method == "GET" {
		return true
	}

	return false
}

func (c *Ctx) IsPost() bool {
	if c.Request().Method == "POST" {
		return true
	}

	return false
}

func (c *Ctx) IsPut() bool {
	if c.Request().Method == "PUT" {
		return true
	}

	return false
}

func (c *Ctx) IsDelete() bool {
	if c.Request().Method == "DELETE" {
		return true
	}

	return false
}

func (c *Ctx) IsPatch() bool {
	if c.Request().Method == "PATCH" {
		return true
	}

	return false
}

func (c *Ctx) IsOptions() bool {
	if c.Request().Method == "OPTIONS" {
		return true
	}

	return false
}

func (c *Ctx) Method() string {
	return c.Request().Method
}

func (c *Ctx) UserAgent() string {
	return c.Request().UserAgent()
}

func (c *Ctx) Host() string {
	return c.Request().Host
}

func (c *Ctx) RemoteAddr() string {
	return c.Request().RemoteAddr
}

func (c *Ctx) IsFromLocal() bool {
	return c.Request().RemoteAddr == ""
}

func (c *Ctx) Header(key Header, def string) string {
	if val := c.Request().Header.Get(string(key)); val != "" {
		return val
	}

	return def
}

func (c *Ctx) SetHeader(key Header, val string) *Ctx {
	c.Response().Header().Set(string(key), val)
	return c
}

func (c *Ctx) GetReqHeaders() map[string][]string {
	return c.Request().Header
}

func (c *Ctx) CurrentRoute() string {
	return c.Request().URL.Path
}
