package restik

import (
	"encoding/json"
	"net/http"
)

type serveReply struct {
	Response interface{} `json:"response,omitempty"`
	Error    *Error      `json:"error,omitempty"`
}

func writeError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	selfErr := FromAnotherError(err)
	reply := serveReply{Error: selfErr}
	writeReply(w, &reply)
	return true
}

func writeResponse(w http.ResponseWriter, resp interface{}) {
	writeReply(w, &serveReply{Response: resp})
}

func writeReply(w http.ResponseWriter, sr *serveReply) {
	b, err := json.Marshal(sr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if sr.Error != nil {
		w.WriteHeader(sr.Error.GetStatus())
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write(b)
}
