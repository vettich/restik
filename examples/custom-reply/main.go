package main

import (
	"net/http"

	"github.com/vettich/restik"
)

type customReply struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

func echo(req *restik.Request) string {
	return req.Vars.String("msg")
}

func main() {
	r := restik.NewRouter()
	r.Get("/echo/{msg}", echo)
	http.Handle("/", r.Handler())
	http.ListenAndServe("0.0.0.0:3303", nil)
}
