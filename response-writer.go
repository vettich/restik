package restik

import (
	"encoding/json"
	"net/http"
)

// ResponseWriter implement rest ResponseWriter from http
type ResponseWriter struct {
	http.ResponseWriter
	jsonHeaderSetted bool
	commonReply      Reply
}

// NewResponseWriter create new ResponseWriter instance
func NewResponseWriter(w http.ResponseWriter, commonReply Reply) ResponseWriter {
	return ResponseWriter{
		ResponseWriter: w,
		commonReply:    commonReply,
	}
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

func (w *ResponseWriter) WriteResponse(resp interface{}) (int, error) {
	if w == nil {
		return 0, nil
	}
	rpl := w.commonReply.New()
	rpl.SetResponse(resp)
	return w.WriteReply(rpl)
}

func (w *ResponseWriter) WriteError(err error) bool {
	if err == nil {
		return false
	}
	selfErr := FromAnotherError(err)
	rpl := w.commonReply.New()
	rpl.SetError(selfErr)
	w.WriteReply(rpl)
	return true
}

func (w *ResponseWriter) WriteReply(rpl Reply) (int, error) {
	b, err := json.Marshal(rpl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := rpl.GetError(); err != nil {
		err := FromAnotherError(err)
		w.WriteHeader(err.GetStatus())
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return w.Write(b)
}
