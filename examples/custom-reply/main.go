package main

import (
	"errors"
	"net/http"

	"github.com/vettich/restik"
)

type customError string

func (err customError) Error() string {
	return string(err)
}

type customReply struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"error,omitempty"`
}

func (cr *customReply) New() restik.Reply {
	return &customReply{}
}

func (cr *customReply) SetResponse(resp interface{}) {
	cr.Data = resp
}

func (cr *customReply) SetError(err error) {
	cerr := customError(err.Error())
	cr.Error = &cerr
}

func (cr *customReply) GetError() error {
	return cr.Error
}

func echo(req *restik.Request) (*string, error) {
	msg, _ := req.Vars.String("msg")
	if msg == "error" {
		return nil, errors.New("error in echo")
	}
	return &msg, nil
}

func main() {
	r := restik.NewRouter()
	r.SetCustomReply(&customReply{})
	r.Get("/echo/{msg}", echo)
	http.Handle("/", r.Handler())
	http.ListenAndServe("0.0.0.0:3303", nil)
}
