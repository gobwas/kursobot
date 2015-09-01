package framework
import (
	"net/http"
	"errors"
	"fmt"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Framework struct {
	handlers map[string]map[string]Handler
}

type Handler func(rw http.ResponseWriter, req *http.Request)

func New() *Framework {
	return &Framework{ handlers:make(map[string]map[string]Handler) }
}

func (self *Framework) Get(path string, handler Handler) {
	self.Use(GET, path, handler)
}

func (self *Framework) Post(path string, handler Handler) {
	self.Use(POST, path, handler)
}

func (self *Framework) Put(path string, handler Handler) {
	self.Use(PUT, path, handler)
}

func (self *Framework) Delete(path string, handler Handler) {
	self.Use(DELETE, path, handler)
}

func (self *Framework) Use(method string, path string, handler Handler)(error) {
	if _, ok := self.handlers[path]; !ok {
		self.handlers[path] = make(map[string]Handler)
	}

	if _, ok := self.handlers[path][method]; ok {
		return errors.New("Handler is already set")
	}

	self.handlers[path][method] = handler

	return nil
}

func abort(rw http.ResponseWriter, code int) {
	rw.WriteHeader(code)
}

func (self *Framework) Listen(port int) {
	for path, desc := range self.handlers {
		http.HandleFunc(path, func(rw http.ResponseWriter, req *http.Request) {
			if handler, ok := desc[req.Method]; ok {
				handler(rw, req)
			} else {
				abort(rw, 405)
			}
		})
	}

	portString := fmt.Sprintf(":%d", port)

	http.ListenAndServe(portString, nil)
}

