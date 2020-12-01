package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vettich/restik"
)

func hello() string {
	return "Hello, world!"
}

func main() {
	r := restik.NewRouter()
	r.Use(&simpleLogger{
		format: "[{METHOD}] {PATH}",
	})
	r.Get("/hello", hello)
	http.Handle("/", r.Handler())
	http.ListenAndServe("0.0.0.0:3303", nil)
}

type simpleLogger struct {
	format string
}

func (mw *simpleLogger) Middleware(next restik.HandlerFunc) restik.HandlerFunc {
	return func(w restik.ResponseWriter, r *restik.Request) {
		logstr := strings.ReplaceAll(mw.format, "{METHOD}", r.Method)
		logstr = strings.ReplaceAll(logstr, "{PATH}", r.URL.RequestURI())
		fmt.Println(logstr)
		next(w, r)
	}
}
