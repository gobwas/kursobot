package app

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type App struct {
	handlers map[string]map[string]http.Handler
}

func New() *App {
	return &App{handlers: make(map[string]map[string]http.Handler)}
}

func (self *App) Get(path string, handler interface{}) error {
	return self.Use(GET, path, handler)
}

func (self *App) Post(path string, handler interface{}) error {
	return self.Use(POST, path, handler)
}

func (self *App) Put(path string, handler interface{}) error {
	return self.Use(PUT, path, handler)
}

func (self *App) Delete(path string, handler interface{}) error {
	return self.Use(DELETE, path, handler)
}

func (self *App) Use(method string, path string, handler interface{}) error {
	var realHandler http.Handler

	switch handler := handler.(type) {
	case http.Handler:
		realHandler = handler
	case func(http.ResponseWriter, *http.Request):
		realHandler = http.HandlerFunc(handler)
	default:
		return errors.New("Handler is expected")
	}

	if _, ok := self.handlers[path]; !ok {
		self.handlers[path] = make(map[string]http.Handler)
	}

	if _, ok := self.handlers[path][method]; ok {
		return errors.New("Handler is already set")
	}

	self.handlers[path][method] = realHandler

	return nil
}

func (self *App) Listen(port int) {
	for path, desc := range self.handlers {
		fmt.Println("registering " + path)

		http.HandleFunc(path, func(rw http.ResponseWriter, req *http.Request) {
			// check for the "/" route
			if req.URL.Path != path {
				http.NotFound(rw, req)
				return
			}

			if handler, ok := desc[req.Method]; ok {
				fmt.Println(fmt.Sprintf("Got %s %s act!", req.Method, path))
				handler.ServeHTTP(rw, req)
			} else {
				http.Error(rw, fmt.Sprintf("%s %s is not implemented", req.Method, path), http.StatusMethodNotAllowed)
			}
		})
	}

	portString := fmt.Sprintf(":%d", port)

	http.ListenAndServe(portString, nil)
}
