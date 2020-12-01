package restik

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(ResponseWriter, *Request)

type Middleware interface {
	Middleware(next HandlerFunc) HandlerFunc
}

type CorsMiddleware struct {
	AllowedMethods []string
	AllowedHeaders []string
	AllowedOrigin  string
}

func (mw *CorsMiddleware) Middleware(next HandlerFunc) HandlerFunc {
	methods := strings.Join(mw.AllowedMethods, ",")
	headers := strings.Join(mw.AllowedHeaders, ",")
	return func(w ResponseWriter, r *Request) {
		log.Println(r.Method, r.URL.Path)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Allow-Origin", mw.AllowedOrigin)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", mw.AllowedOrigin)
		next(w, r)
	}
}
