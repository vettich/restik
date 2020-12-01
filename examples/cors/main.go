package main

import (
	"net/http"

	"github.com/vettich/restik"
)

func hello() string {
	return "Hello, world!"
}

func main() {
	r := restik.NewRouter()
	r.Use(&restik.CorsMiddleware{
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Origin"},
		AllowedOrigin:  "*",
	})
	r.Get("/hello", hello)
	http.Handle("/", r.Handler())
	http.ListenAndServe("0.0.0.0:3303", nil)
}
