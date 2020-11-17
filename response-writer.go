package restik

import (
	"encoding/json"
	"net/http"
)

// ResponseWriter implement rest ResponseWriter from http
type ResponseWriter struct {
	http.ResponseWriter
	jsonHeaderSetted bool
}

// NewResponseWriter create new ResponseWriter instance
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{ResponseWriter: w}
}

// WriteJSON write src to response with encoding json
func (w *ResponseWriter) WriteJSON(src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if !w.jsonHeaderSetted {
		w.Header().Set("Content-Type", "application/json")
		w.jsonHeaderSetted = true
	}
	_, err = w.Write(b)
	return err
}
