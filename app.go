package main

import (
	"github.com/gobwas/kursobot/framework"
	"net/http"
)

func main() {
	app := framework.New()

	// register handler
	app.Get("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("Hi!"))
	})

	app.Listen(8080)
}