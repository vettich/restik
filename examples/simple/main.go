package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/vettich/restik"
)

func echo(req *restik.Request) string {
	return req.Vars.String("msg")
}

type helloArg struct {
	Value string `json:"value"`
}

func hello(arg helloArg) (string, error) {
	if arg.Value == "" {
		return "", errors.New("value is empty")
	}
	return "hello, " + arg.Value, nil
}

func httpFn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("usuallyHandlerFn runned"))
}

func restFn(w restik.ResponseWriter, r *restik.Request) {
	w.WriteJSON(r.Vars.String("you"))
}

func loggg() {
	log.Println("logging")
}

func main() {
	r := restik.NewRouter()
	r.Get("/echo/{msg}", echo)
	r.Post("/hello", hello)
	r.Get("/log", loggg)
	r.Get("/http", httpFn)
	r.Get("/src/{you}", restFn)
	http.Handle("/", r.Handler())
	http.ListenAndServe("0.0.0.0:3303", nil)
}
