package main

import (
	"errors"
	"log"
	"net/http"
	"restik"
)

type echoArg struct {
	Value string `json:"value"`
}

func echo(arg echoArg) (string, error) {
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
	router := restik.NewRouter()
	router.Add(
		restik.Post("/echo", echo),
		restik.Get("/log", loggg),
		restik.Get("/http", httpFn),
		restik.Get("/src/{you}", restFn),
	)
	http.Handle("/", router.Handler())
	http.ListenAndServe("0.0.0.0:3303", nil)
}
